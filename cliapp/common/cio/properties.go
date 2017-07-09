package cio

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/varutil"
)

// ReadProperties read properties from app.Input
func ReadProperties(baseKey string, in app.Input, out app.Output, def []*config.Property, data, defaultData map[string]string) (isChanged bool, err error) {
	var (
		ok           bool
		defaultValue string
		input        string
	)
	for _, property := range def {
		if _, ok = data[property.Key]; ok {
			continue
		}
		if defaultValue, ok = defaultData[property.Key]; !ok {
			switch strings.ToLower(property.Type) {
			case "numeric":
				defaultValue = varutil.RandString(property.Max, varutil.NumericBytes)
			case "alpha":
				defaultValue = varutil.RandString(property.Max, varutil.AlphaBytes)
			case "alnum":
				defaultValue = varutil.RandString(property.Max, varutil.AlphaNumericBytes)
			case "strong":
				defaultValue = varutil.RandString(property.Max, varutil.StrongBytes)
			case "line":
				defaultValue = ""
			default:
				return isChanged, fmt.Errorf("wrong property type %s (for property %s)", property.Type, property.Key)
			}
		}
		for {
			prompt := property.Prompt
			if prompt == "" {
				prompt = "Insert property"
			}
			if defaultValue != "" {
				out.Printf("(%s) %s [%s]: ", property.Key, prompt, defaultValue)
			} else {
				out.Printf("(%s) %s: ", property.Key, prompt)
			}
			if input, err = in.ReadLine(); err != nil && err != io.EOF {
				return isChanged, err
			}
			if input == "" && defaultValue != "" {
				isChanged = true
				data[property.Key] = defaultValue
				break
			}
			if len(input) < property.Min {
				out.Printf("Value is too short. Minimum length of the value is %d characters.\n", property.Min)
				continue
			}
			if len(input) > property.Max {
				out.Printf("Value is too long. Maximum length of the value is %d characters.\n", property.Max)
				continue
			}
			if property.Type == "numeric" && !numericReg.MatchString(input) {
				out.Printf("Require numeric value ('%s' is incorrect)\n", input)
				continue
			}
			if property.Type == "alpha" && !alphaReg.MatchString(input) {
				out.Printf("Require alpha value ('%s' is incorrect)\n", input)
				continue
			}
			if property.Type == "alnum" && !alnumReg.MatchString(input) {
				out.Printf("Require alpha-numeric value ('%s' is incorrect)\n", input)
				continue
			}
			if property.Pattern != nil && !property.Pattern.MatchString(input) {
				out.Printf("Invalid value. Value must be matched by '%s' regexp\n", property.Pattern.String())
				continue
			}
			isChanged = true
			data[baseKey+property.Key] = input
			break
		}
	}
	return isChanged, nil
}

// ReadPropertiesArray read properties array
func ReadPropertiesArray(baseKey string, in app.Input, out app.Output, properties []*config.Property, data map[string]string) (isChanged bool, err error) {
	var (
		line          string
		isCollChanged bool
		i             int
	)
loop:
	for {
		out.Printf("Do you want add a element? (y/n): ")
		if line, err = in.ReadLine(); err != nil {
			return isChanged, err
		}
		switch strings.ToLower(line) {
		case "y":
			subKey := baseKey + strconv.Itoa(i) + "."
			if isCollChanged, err = ReadProperties(subKey, in, out, properties, data, map[string]string{}); err != nil {
				return isChanged, err
			}
			isChanged = isCollChanged || isChanged
		case "n":
			break loop
		default:
			continue
		}
		i++
	}
	return isChanged, nil
}

// ReadPropertiesMap read properties map
func ReadPropertiesMap(baseKey string, in app.Input, out app.Output, properties []*config.Property, data map[string]string) (isChanged bool, err error) {
	var (
		mapkey        string
		isCollChanged bool
		i             int
	)
	for {
		out.Printf("Do you want add a element? Insert a key or empty value to break: ")
		if mapkey, err = in.ReadLine(); err != nil {
			if err == io.EOF {
				return isChanged, nil
			}
			return isChanged, err
		}
		if mapkey == "" {
			return isChanged, nil
		}
		subKey := baseKey + mapkey + "."
		if isCollChanged, err = ReadProperties(subKey, in, out, properties, data, map[string]string{}); err != nil {
			return isChanged, err
		}
		isChanged = isCollChanged || isChanged
		i++
	}
}
