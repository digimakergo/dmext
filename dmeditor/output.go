package dmeditor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query/querier"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

// type ContentParams struct {
// 	Mode        string `json:"mode"`
// 	ContentType string `json:"contenttype"`
// 	ID          int    `json:"id"`
// }

func ProceedData(ctx context.Context, jsonField fieldtypes.Json) (string, error) {
	dmeditorServer := viper.GetString("general.dmeditor_server_url")
	if dmeditorServer == "" {
		return "", errors.New("No dmeditor server configured")
	}

	//send request with json
	jsonBytes := jsonField.Content
	bodyResult := ""
	if jsonBytes != nil {
		resp, err := http.Post(dmeditorServer, "application/json", bytes.NewBuffer(jsonBytes))

		if err != nil {
			return "", fmt.Errorf("Proceeding json error: %v", err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		bodyResult = string(body)

		// //wash dmcontent:<json>
		// dmeditorWash := viper.GetBoolean("general.dmeditor_wash")
		// if dmeditorWash {
		// 	reg, _ := regexp.Compile("<!--dmcontent:(.*?)-->")
		// 	resultList := reg.FindAll(body, -1)
		// 	for _, result := range resultList {
		// 		resultStr := string(result)
		// 		contentParamsStr := strings.TrimRight(strings.TrimLeft(resultStr, "<!--dmcontent:"), "-->")

		// 		params := ContentParams{}
		// 		err := json.Unmarshal([]byte(contentParamsStr), &params)

		// 		vars := map[string]interface{}{}
		// 		content, _ := query.FetchByID(ctx, params.ID)
		// 		vars["content"] = content
		// 		vars["viewmode"] = "editor_block" //params.Mode
		// 		str, err := sitekit.OutputString(vars, "content_view", sitekit.RequestInfo{Context: ctx, Site: "dmdemo"})

		// 		bodyResult = strings.ReplaceAll(bodyResult, resultStr, str)

		// 	}
		// 	return bodyResult, nil
		// }
	}

	//proceed with html
	return bodyResult, nil
}

type DMEditorOutputer struct {
}

func (d DMEditorOutputer) Output(ctx context.Context, querier querier.Querier, value interface{}) interface{} {
	result, err := ProceedData(ctx, value.(fieldtypes.Json))
	if err != nil {
		//todo: return err
		fmt.Println(err.Error())
		log.Error(err.Error(), "", ctx)
		return ""
	}
	return result
}

func init() {
	outputer := DMEditorOutputer{}
	fieldtypes.RegisterJSONOutputer("dmeditor", outputer)
	log.Info("Registering dmeditor json outputer")
}
