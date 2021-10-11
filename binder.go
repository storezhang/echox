package echox

import (
	`bytes`
	`encoding`
	`encoding/json`
	`encoding/xml`
	`errors`
	`net/http`
	`reflect`
	`strconv`
	`strings`

	`github.com/labstack/echo/v4`
	`github.com/mcuadros/go-defaults`
	`github.com/vmihailenco/msgpack/v5`
	`google.golang.org/protobuf/proto`
)

type binder struct {
	tagParam  string
	tagQuery  string
	tagForm   string
	tagHeader string
}

func (b *binder) Bind(value interface{}, ctx echo.Context) (err error) {
	if err = b.params(ctx, value); nil != err {
		return
	}
	if err = b.headers(ctx, value); nil != err {
		return
	}

	if http.MethodGet == ctx.Request().Method || http.MethodDelete == ctx.Request().Method {
		if err = b.queries(ctx, value); nil != err {
			return
		}
	}

	// 只有在Content-Type设置值后才绑定Body
	contentType := ctx.Request().Header.Get(HeaderContentType)
	if "" == contentType {
		return
	}
	if err = b.body(ctx, contentType, value); nil != err {
		return
	}

	if reflect.Ptr == reflect.ValueOf(value).Kind() {
		defaults.SetDefaults(value)
	} else {
		defaults.SetDefaults(&value)
	}

	return
}

func (b *binder) params(ctx echo.Context, value interface{}) (err error) {
	names := ctx.ParamNames()
	values := ctx.ParamValues()
	params := map[string][]string{}
	for index, name := range names {
		params[name] = []string{values[index]}
	}

	if err = b.bindData(value, params, b.tagParam); nil != err {
		err = echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return
}

func (b *binder) queries(ctx echo.Context, value interface{}) (err error) {
	if err = b.bindData(value, ctx.QueryParams(), b.tagQuery); nil != err {
		err = echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return
}

func (b *binder) headers(ctx echo.Context, i interface{}) (err error) {
	if err = b.bindData(i, ctx.Request().Header, b.tagHeader); nil != err {
		err = echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return
}

func (b *binder) body(ctx echo.Context, contentType string, value interface{}) (err error) {
	req := ctx.Request()
	if req.ContentLength == 0 {
		return
	}

	switch {
	case strings.HasPrefix(contentType, MIMEApplicationJSON):
		err = json.NewDecoder(req.Body).Decode(value)
	case strings.HasPrefix(contentType, MIMEApplicationXML), strings.HasPrefix(contentType, MIMETextXML):
		err = xml.NewDecoder(req.Body).Decode(value)
	case strings.HasPrefix(contentType, MIMEApplicationProtobuf):
		buf := new(bytes.Buffer)
		if _, err = buf.ReadFrom(req.Body); nil != err {
			return
		}
		err = proto.Unmarshal(buf.Bytes(), value.(proto.Message))
	case strings.HasPrefix(contentType, MIMEApplicationMsgpack):
		err = msgpack.NewDecoder(req.Body).Decode(value)
	case strings.HasPrefix(contentType, MIMEApplicationForm), strings.HasPrefix(contentType, MIMEMultipartForm):
		var params map[string][]string
		if params, err = ctx.FormParams(); nil != err {
			return
		}
		err = b.bindData(value, params, b.tagForm)
	}

	return
}

func (b *binder) bindData(destination interface{}, data map[string][]string, tag string) (err error) {
	if nil == destination || 0 == len(data) {
		return
	}
	typ := reflect.TypeOf(destination).Elem()
	val := reflect.ValueOf(destination).Elem()

	if typ.Kind() == reflect.Map {
		for k, v := range data {
			val.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
		}

		return
	}

	if typ.Kind() != reflect.Struct {
		if tag == b.tagParam || tag == b.tagQuery || tag == b.tagHeader {
			// incompatible type, data is probably to be found in the body
			return nil
		}
		return errors.New("binding element must be a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if typeField.Anonymous {
			if structField.Kind() == reflect.Ptr {
				structField = structField.Elem()
			}
		}
		if !structField.CanSet() {
			continue
		}
		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(tag)
		if typeField.Anonymous && structField.Kind() == reflect.Struct && inputFieldName != "" {
			// if anonymous struct with query/param/form tags, report an error
			return errors.New("query/param/form tags are not allowed with anonymous struct field")
		}

		if inputFieldName == "" {
			// If tag is nil, we inspect if the field is a not BindUnmarshaler struct and try to bind data into it (might contains fields with tags).
			// structs that implement BindUnmarshaler are binded only when they have explicit tag
			if _, ok := structField.Addr().Interface().(echo.BindUnmarshaler); !ok && structFieldKind == reflect.Struct {
				if err := b.bindData(structField.Addr().Interface(), data, tag); err != nil {
					return err
				}
			}
			// does not have explicit tag and is not an ordinary struct - so move to next field
			continue
		}

		inputValue, exists := data[inputFieldName]
		if !exists {
			// Go json.Unmarshal supports case insensitive binding.  However the
			// url params are bound case sensitive which is inconsistent.  To
			// fix this we must check all of the map values in a
			// case-insensitive search.
			for k, v := range data {
				if strings.EqualFold(k, inputFieldName) {
					inputValue = v
					exists = true
					break
				}
			}
		}

		if !exists {
			continue
		}

		// Call this first, in case we're dealing with an alias to an array type
		if ok, err := unmarshalField(typeField.Type.Kind(), inputValue[0], structField); ok {
			if err != nil {
				return err
			}
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for j := 0; j < numElems; j++ {
				if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
			return err

		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) (err error) {
	var ok bool
	if ok, err = unmarshalField(valueKind, val, structField); ok {
		return
	}

	switch valueKind {
	case reflect.Ptr:
		err = setWithProperType(structField.Elem().Kind(), val, structField.Elem())
	case reflect.Int:
		err = setIntField(val, 0, structField)
	case reflect.Int8:
		err = setIntField(val, 8, structField)
	case reflect.Int16:
		err = setIntField(val, 16, structField)
	case reflect.Int32:
		err = setIntField(val, 32, structField)
	case reflect.Int64:
		err = setIntField(val, 64, structField)
	case reflect.Uint:
		err = setUintField(val, 0, structField)
	case reflect.Uint8:
		err = setUintField(val, 8, structField)
	case reflect.Uint16:
		err = setUintField(val, 16, structField)
	case reflect.Uint32:
		err = setUintField(val, 32, structField)
	case reflect.Uint64:
		err = setUintField(val, 64, structField)
	case reflect.Bool:
		err = setBoolField(val, structField)
	case reflect.Float32:
		err = setFloatField(val, 32, structField)
	case reflect.Float64:
		err = setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		err = errors.New("unknown type")
	}

	return
}

func unmarshalField(valueKind reflect.Kind, val string, field reflect.Value) (bool, error) {
	switch valueKind {
	case reflect.Ptr:
		return unmarshalFieldPtr(val, field)
	default:
		return unmarshalFieldNonPtr(val, field)
	}
}

func unmarshalFieldNonPtr(value string, field reflect.Value) (bool, error) {
	fieldIValue := field.Addr().Interface()
	if unmarshaler, ok := fieldIValue.(echo.BindUnmarshaler); ok {
		return true, unmarshaler.UnmarshalParam(value)
	}
	if unmarshaler, ok := fieldIValue.(encoding.TextUnmarshaler); ok {
		return true, unmarshaler.UnmarshalText([]byte(value))
	}

	return false, nil
}

func unmarshalFieldPtr(value string, field reflect.Value) (bool, error) {
	if field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
	}

	return unmarshalFieldNonPtr(value, field.Elem())
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if "" == value {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if "" == value {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if "" == value {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if "" == value {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}
