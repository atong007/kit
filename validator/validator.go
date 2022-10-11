package validator

import (
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/zh"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Validate value validation
var val *validator.Validate
var trans ut.Translator

func init() {
	zhLocales := zh.New()
	uni := ut.New(zhLocales, zhLocales)
	trans, _ = uni.GetTranslator("zh")
	val = validator.New()
	val.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	err := val.RegisterValidation("mobile", validateMobile)
	if err != nil {
		log.Fatalf("error registering valitor: %v", err)
	}

	registerMobileErrTrans()

	// 验证器注册翻译器
	err = zhTrans.RegisterDefaultTranslations(val, trans)
	if err != nil {
		log.Fatalf("error registering default translation for validator: %w", err)
	}
}

func Check(s interface{}) error {
	return val.Struct(s)
}

func registerMobileErrTrans() {
	_ = val.RegisterTranslation("mobile", trans,
		func(ut ut.Translator) error {
			return ut.Add("mobile", "手机号格式错误", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			fld, _ := ut.T(fe.Field())
			t, err := ut.T(fe.Tag(), fld)
			if err != nil {
				return fe.(error).Error()
			}
			return t
		},
	)
}

// Translate translate en to zh
func Translate(valErr error) string {
	if valErr, ok := valErr.(validator.ValidationErrors); ok {
		var errList []string
		for _, e := range valErr {
			// can translate each error one at a time.
			errList = append(errList, e.Translate(trans))
		}
		return strings.Join(errList, "|")
	}
	return valErr.Error()
}

// 验证手机号
func validateMobile(fl validator.FieldLevel) bool {
	patterns := map[string]string{
		"86":  `^1[345789]\d{9}$`,             // 中文（中华人民共和国）
		"886": `^9\d{8}$`,                     // 中文（台湾）
		"213": `^(5|6|7)\d{8}$`,               // 阿拉伯文（阿尔及利亚）
		"963": `^9\d{8}$`,                     // 阿拉伯文（叙利亚）
		"966": `^5\d{8}$`,                     // 阿拉伯语（沙特阿拉伯）
		"1":   `^[2-9]\d{2}[2-9](?!11)\d{6}$`, // 英语（美国）
	}

	countryCode := reflect.Indirect(fl.Top()).FieldByName("CountryCode")
	if !countryCode.IsValid() {
		return false
	}
	pattern, ok := patterns[countryCode.String()]
	if !ok {
		pattern = patterns["86"]
	}

	matched, err := regexp.MatchString(pattern, fl.Field().String())
	if matched == false || err != nil {
		return false
	}
	return true
}
