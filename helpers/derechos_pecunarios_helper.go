package helpers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

// FiltrarDerechosPecuniarios ...
// @Title FiltrarDerechosPecuniarios
// @Description Consulta los parametros y filtra los conceptos de derechos pecuniarios a partir del Id de la vigencia
func FiltrarDerechosPecuniarios(vigenciaId string) ([]interface{}, error) {
	var parametros map[string]interface{}
	var conceptos []interface{}

	errorConceptos := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?limit=0&query=PeriodoId__Id:"+vigenciaId, &parametros)
	if errorConceptos == nil {
		if parametros["Data"] != nil && fmt.Sprintf("%v", parametros["Data"]) != "[map[]]" {
			conceptos = parametros["Data"].([]interface{})
			if fmt.Sprintf("%v", conceptos[0]) != "map[]" {
				conceptosFiltrados := conceptos[:0]
				for _, concepto := range conceptos {
					TipoParametro := concepto.(map[string]interface{})["ParametroId"].(map[string]interface{})["TipoParametroId"].(map[string]interface{})["Id"].(float64)
					if TipoParametro == 2 && concepto.(map[string]interface{})["Activo"] == true { //id para derechos_pecuniarios
						conceptosFiltrados = append(conceptosFiltrados, concepto)
					}
				}
				conceptos = conceptosFiltrados
			}
		}
	}
	return conceptos, errorConceptos
}
