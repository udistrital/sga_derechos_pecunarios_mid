package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid_derechos_pecunarios/helpers"
	"github.com/udistrital/sga_mid_derechos_pecunarios/models"
	"github.com/udistrital/utils_oas/request"
)

func PostConcepto(data []byte) (interface{}, error) {
	var ConceptoFactor map[string]interface{}
	var AuxConceptoPost map[string]interface{}
	var ConceptoPost map[string]interface{}
	var IdConcepto interface{}
	var NumFactor interface{}
	var Vigencia interface{}
	var ValorJson map[string]interface{}
	var exito bool = true
	var respuesta interface{}

	//Se guarda el json que se pasa por parametro
	if err := json.Unmarshal(data, &ConceptoFactor); err == nil {
		Concepto := ConceptoFactor["Concepto"]
		errConcepto := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro", "POST", &AuxConceptoPost, Concepto)
		ConceptoPost = AuxConceptoPost["Data"].(map[string]interface{})
		IdConcepto = ConceptoPost["Id"]

		if errConcepto == nil && fmt.Sprintf("%v", ConceptoPost["System"]) != "map[]" && ConceptoPost["Id"] != nil {
			if ConceptoPost["Status"] == 400 {
				logs.Error(errConcepto)
				exito = false
			} 
		} else {
			logs.Error(errConcepto)
			exito = false
		}

		Vigencia = ConceptoFactor["Vigencia"] //Valor del id de la vigencia (periodo)
		NumFactor = ConceptoFactor["Factor"]  //Valor que trae el numero del factor y el salario minimo

		ValorFactor := fmt.Sprintf("%.3f", NumFactor.(map[string]interface{})["Valor"].(map[string]interface{})["NumFactor"])
		Valor := "{\n    \"NumFactor\": " + ValorFactor + " \n}"

		Factor := map[string]interface{}{
			"ParametroId": map[string]interface{}{"Id": IdConcepto.(float64)},
			"PeriodoId":   map[string]interface{}{"Id": Vigencia.(map[string]interface{})["Id"].(float64)},
			"Valor":       Valor,
			"Activo":      true,
		}

		var AuxFactor map[string]interface{}
		var FactorPost map[string]interface{}

		errFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo", "POST", &AuxFactor, Factor)
		FactorPost = AuxFactor["Data"].(map[string]interface{})
		if errFactor == nil && fmt.Sprintf("%v", FactorPost["System"]) != "map[]" && FactorPost["Id"] != nil {
			if FactorPost["Status"] != 400 {
				//JSON que retorna al agregar el concepto y el factor
				ValorString := FactorPost["Valor"].(string)
				if err := json.Unmarshal([]byte(ValorString), &ValorJson); err == nil {
					respuesta = map[string]interface{}{
						"Concepto": map[string]interface{}{
							"Id":                IdConcepto.(float64),
							"Nombre":            ConceptoPost["Nombre"],
							"CodigoAbreviacion": ConceptoPost["CodigoAbreviacion"],
							"Activo":            ConceptoPost["Activo"],
						},
						"Factor": map[string]interface{}{
							"Id":    FactorPost["Id"],
							"Valor": ValorJson["NumFactor"],
						},
					}
				}
			} else {
				var resultado2 map[string]interface{}
				request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("ParametroService")+"parametro/%.f", ConceptoPost["Id"]), "DELETE", &resultado2, nil)
				logs.Error(errFactor)
				exito = false
			}
		} else {
			logs.Error(errFactor)
			exito = false
		}
	} else {
		exito = false
	}
	if exito {
		return respuesta, nil
	}
	return nil, errors.New("error del servicio PostConcepto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func PutConcepto(idConcepto string, data []byte) (interface{}, error) {
	var ConceptoFactor map[string]interface{}
	var AuxConceptoPut map[string]interface{}
	var AuxFactorPut map[string]interface{}
	var ConceptoPut map[string]interface{}
	var Parametro map[string]interface{}
	var respuesta interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId__Id:"+idConcepto, &Parametro); err == nil {
		DataAux := Parametro["Data"].([]interface{})[0]
		Data := DataAux.(map[string]interface{})

		if fmt.Sprintf("%v", Data) != "map[]" {
			ConceptoPut = Data["ParametroId"].(map[string]interface{})
			if err := json.Unmarshal(data, &ConceptoFactor); err == nil {
				if fmt.Sprintf("%v", ConceptoFactor) != "map[]" {
					Factor := ConceptoFactor["Factor"].(map[string]interface{})
					FactorValor := fmt.Sprintf("%.3f", Factor["Valor"].(map[string]interface{})["NumFactor"].(float64))
					Data["Valor"] = "{ \"NumFactor\": " + FactorValor + " }"
					errFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo/"+fmt.Sprintf("%.f", Data["Id"].(float64)), "PUT", &AuxFactorPut, Data)
					if errFactor != nil {
						logs.Error(errFactor)
					}
					Concepto := ConceptoFactor["Concepto"].(map[string]interface{})
					ConceptoPut["Nombre"] = Concepto["Nombre"]
					ConceptoPut["CodigoAbreviacion"] = Concepto["CodigoAbreviacion"]
					errPut := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/"+idConcepto, "PUT", &AuxConceptoPut, ConceptoPut)
					if errPut != nil {
						logs.Error(errPut)
					} else {
						respuesta = map[string]interface{}{
							"Concepto": AuxConceptoPut,
							"Factor":   AuxFactorPut,
						}
						return respuesta, nil
					}
				} else {
					return nil, errors.New("error del servicio PutConcepto: Body is empty")
				}
			} else {
				logs.Error(err)
				return nil, errors.New("error del servicio PutConcepto: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		} else {
			return nil, errors.New("error del servicio PutConcepto: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	} else {
		logs.Error(err)
		return nil, errors.New("error del servicio PutConcepto: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
	return nil, errors.New("error del servicio PutConcepto: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func DeleteConcepto(idConcepto string) (interface{}, error) {
	var Parametro map[string]interface{}
	var AuxFactorPut map[string]interface{}
	var AuxConceptoPut map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId__Id:"+idConcepto, &Parametro); err == nil {
		DataAux := Parametro["Data"].([]interface{})[0]
		Data := DataAux.(map[string]interface{})
		Concepto := Data["ParametroId"].(map[string]interface{})
		Data["Activo"] = false
		Concepto["Activo"] = false

		errFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo/"+fmt.Sprintf("%.f", Data["Id"].(float64)), "PUT", &AuxFactorPut, Data)
		if errFactor == nil {
			errConcepto := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/"+idConcepto, "PUT", &AuxConceptoPut, Concepto)
			if errConcepto == nil {

				response := map[string]interface{}{
					"Concepto": AuxConceptoPut,
					"Factor":   AuxFactorPut,
				}
				return response, nil
			} else {
				logs.Error(errConcepto)
				return nil, errors.New("error del servicio DeleteConcepto: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		} else {
			logs.Error(errFactor)
			return nil, errors.New("error del servicio DeleteConcepto: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	} else {
		logs.Error(err)
		return nil, errors.New("error del servicio DeleteConcepto: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
}

func GetDerechosPecuniariosPorVigencia(idVigencia string) (interface{}, error) {
	var conceptos []interface{}
	var err error
	conceptos, err = helpers.FiltrarDerechosPecuniarios(idVigencia)

	if err == nil {
		if conceptos != nil {
			return conceptos, nil
		} else {
			return nil, errors.New("error del servicio GetDerechosPecuniariosPorVigencia: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	} else {
		logs.Error(err)
		return nil, errors.New("error del servicio GetDerechosPecuniariosPorVigencia: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
}

func PostClonarConceptos(data []byte) (interface{}, error) {
	var vigencias map[string]interface{}
	var conceptos []interface{}
	var NuevoConceptoPost map[string]interface{}
	var NuevoFactorPost map[string]interface{}
	var errorConceptos error
	var errorGetAll bool

	if errorVigencias := json.Unmarshal(data, &vigencias); errorVigencias == nil {
		vigenciaAnterior := vigencias["VigenciaAnterior"].(float64)
		vigenciaActual := vigencias["VigenciaActual"].(float64)
		conceptos, errorConceptos = helpers.FiltrarDerechosPecuniarios(fmt.Sprintf("%.f", vigenciaAnterior))
		if errorConceptos == nil {
			for _, concepto := range conceptos {
				OldConcepto := concepto.(map[string]interface{})["ParametroId"].(map[string]interface{})
				TipoParametroId := OldConcepto["TipoParametroId"].(map[string]interface{})["Id"].(float64)
				NuevoConcepto := map[string]interface{}{
					"Nombre":            OldConcepto["Nombre"],
					"Descripcion":       OldConcepto["Descripcion"],
					"CodigoAbreviacion": OldConcepto["CodigoAbreviacion"],
					"NumeroOrden":       OldConcepto["NumeroOrden"],
					"Activo":            OldConcepto["Activo"],
					"TipoParametroId":   map[string]interface{}{"Id": TipoParametroId},
				}
				errNuevoConcepto := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro", "POST", &NuevoConceptoPost, NuevoConcepto)
				if errNuevoConcepto == nil {
					OldFactor := concepto.(map[string]interface{})
					NuevoFactor := map[string]interface{}{
						"Valor":       OldFactor["Valor"],
						"Activo":      OldFactor["Activo"],
						"ParametroId": map[string]interface{}{"Id": NuevoConceptoPost["Data"].(map[string]interface{})["Id"]},
						"PeriodoId":   map[string]interface{}{"Id": vigenciaActual},
					}
					errNuevoFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo", "POST", &NuevoFactorPost, NuevoFactor)
					if errNuevoFactor != nil {
						var resDelete string
						errorGetAll = true
						logs.Error(errNuevoFactor)
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("ParametroService")+"parametro/%.f", NuevoConceptoPost["Id"]), "DELETE", &resDelete, nil)
					}
				} else {
					errorGetAll = true
					logs.Error(errNuevoConcepto)
				}
			}
		} else {
			errorGetAll = true
			logs.Error(errorConceptos)
		}
	} else {
		errorGetAll = true
	}

	if !errorGetAll {
		return NuevoFactorPost, nil
	}
	return nil, errors.New("error del servicio PostClonarConceptos: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func PutCostoConcepto(data []byte) (interface{}, error) {
	var ConceptoCostoAux []map[string]interface{}
	var Factor map[string]interface{}
	var FactorPut map[string]interface{}
	var FactorAux map[string]interface{}
	var errorGetAll bool

	//Guarda el arreglo de objetos  de los conceptos que se traen del cliente
	if err := json.Unmarshal(data, &ConceptoCostoAux); err == nil {
		//Recorre cada concepto para poder guardar el costo
		if fmt.Sprintf("%v", ConceptoCostoAux) != "[map[]]" && fmt.Sprintf("%v", ConceptoCostoAux) != "[]" {
			for _, conceptoTemp := range ConceptoCostoAux {
				idFactor := fmt.Sprintf("%.f", conceptoTemp["FactorId"].(float64))
				// Consulta el factor que esta relacionado con el valor del concepto
				errFactor := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo/"+idFactor, &FactorAux)
				if errFactor == nil {
					if FactorAux != nil {
						Factor = FactorAux["Data"].(map[string]interface{})
						FactorValor := fmt.Sprintf("%.3f", conceptoTemp["Factor"].(float64))
						CostoValor := fmt.Sprintf("%.f", conceptoTemp["Costo"].(float64))
						Valor := "{\n    \"NumFactor\": " + FactorValor + ",\n \"Costo\": " + CostoValor + "\n}"
						Factor["Valor"] = Valor
						errPut := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo/"+idFactor, "PUT", &FactorPut, Factor)
						if errPut == nil {
							if FactorPut == nil {
								errorGetAll = true
							}
						} else {
							errorGetAll = true
							logs.Error(errPut)
						}
					} else {
						errorGetAll = true
						logs.Error(errFactor)
					}
				} else {
					errorGetAll = true
					logs.Error(errFactor)
				}
			}
		} else {
			errorGetAll = true
		}
	} else {
		errorGetAll = true
	}

	if !errorGetAll {
		return FactorPut, nil
	}
	return nil, errors.New("error del servicio PutCostoConcepto: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func PostGenerarDerechoPecuniarioEstudiante(data []byte) (interface{}, error) {
	var SolicitudDerechoPecuniario map[string]interface{}
	var TipoParametro string
	var Derecho map[string]interface{}
	var Codigo []interface{}
	var Valor map[string]interface{}
	var NuevoRecibo map[string]interface{}
	var complementario map[string]interface{}
	var errorGetAll bool

	if err := json.Unmarshal(data, &SolicitudDerechoPecuniario); err == nil {
		if fmt.Sprintf("%v", SolicitudDerechoPecuniario) != "map[]" {
			objTransaccion := map[string]interface{}{
				"codigo":              "-------",
				"nombre":              SolicitudDerechoPecuniario["Nombre"].(string),
				"apellido":            SolicitudDerechoPecuniario["Apellido"].(string),
				"correo":              SolicitudDerechoPecuniario["Correo"].(string),
				"proyecto":            SolicitudDerechoPecuniario["ProgramaAcademicoId"].(string),
				"tiporecibo":          0,
				"concepto":            "-------",
				"valorordinario":      0,
				"valorextraordinario": 0,
				"cuota":               1,
				"fechaordinario":      SolicitudDerechoPecuniario["FechaPago"].(string),
				"fechaextraordinario": SolicitudDerechoPecuniario["FechaPago"].(string),
				"aniopago":            SolicitudDerechoPecuniario["Year"].(float64),
				"perpago":             SolicitudDerechoPecuniario["Periodo"].(float64),
			}

			paramId := fmt.Sprintf("%.f", SolicitudDerechoPecuniario["DerechoPecuniarioId"].(float64))
			terceroId := fmt.Sprintf("%.f", SolicitudDerechoPecuniario["Id"].(float64))
			errParam := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId.Id:"+paramId, &Derecho)
			if errParam == nil && fmt.Sprintf("%v", Derecho["Data"].([]interface{})[0]) != "map[]" {

				errCodigo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=InfoComplementariaId.Id:93,TerceroId.Id:"+terceroId, &Codigo)
				if errCodigo == nil && fmt.Sprintf("%v", Codigo) != "map[]" {
					objTransaccion["codigo"] = Codigo[0].(map[string]interface{})["Dato"]

					Dato := Derecho["Data"].([]interface{})[0]
					if errJson := json.Unmarshal([]byte(Dato.(map[string]interface{})["Valor"].(string)), &Valor); errJson == nil {
						objTransaccion["valorordinario"] = Valor["Costo"].(float64)
						objTransaccion["valorextraordinario"] = Valor["Costo"].(float64)

						TipoParametro = fmt.Sprintf("%v", Dato.(map[string]interface{})["ParametroId"].(map[string]interface{})["CodigoAbreviacion"])
						// Pendiente SISTEMATICACION, MULTAS BIBLIOTECA y FOTOCOPIAS
						switch TipoParametro {
						case "40":
							objTransaccion["tiporecibo"] = 5
							objTransaccion["concepto"] = "CERTIFICADO DE NOTAS"
						case "50":
							objTransaccion["tiporecibo"] = 8
							objTransaccion["concepto"] = "DERECHOS DE GRADO"
						case "51":
							objTransaccion["tiporecibo"] = 9
							objTransaccion["concepto"] = "DUPLICADO DEL DIPLOMA DE GRADO"
						case "44":
							objTransaccion["tiporecibo"] = 10
							objTransaccion["concepto"] = "DUPLICADO DEL CARNET ESTUDIANTIL"
						case "31":
							objTransaccion["tiporecibo"] = 13
							objTransaccion["concepto"] = "CURSOS VACIONALES"
						case "41":
							objTransaccion["tiporecibo"] = 6
							objTransaccion["concepto"] = "CONSTANCIAS DE ESTUDIO"
						case "49":
							objTransaccion["tiporecibo"] = 17
							objTransaccion["concepto"] = "COPIA ACTA DE GRADO"
						case "42":
							objTransaccion["tiporecibo"] = 18
							objTransaccion["concepto"] = "CARNET ESTUDIANTIL"
						}

						SolicitudRecibo := objTransaccion

						reciboSolicitud := httplib.Post("http://" + beego.AppConfig.String("GenerarReciboJbpmService") + "recibos_pago_proxy")
						reciboSolicitud.Header("Accept", "application/json")
						reciboSolicitud.Header("Content-Type", "application/json")
						reciboSolicitud.JSONBody(SolicitudRecibo)

						if errRecibo := reciboSolicitud.ToJSON(&NuevoRecibo); errRecibo == nil {
							derechoPecuniarioSolicitado := map[string]interface{}{
								"TerceroId": map[string]interface{}{
									"Id": SolicitudDerechoPecuniario["Id"].(float64),
								},
								"InfoComplementariaId": map[string]interface{}{
									"Id": 307,
								},
								"Activo": true,
								"Dato":   `{"Recibo":` + `"` + fmt.Sprintf("%v/%v", NuevoRecibo["creaTransaccionResponse"].(map[string]interface{})["secuencia"], NuevoRecibo["creaTransaccionResponse"].(map[string]interface{})["anio"]) + `", ` + `"CodigoAsociado": "` + SolicitudDerechoPecuniario["CodigoEstudiante"].(string) + `", "SolicitudId":""}`,
							}

							errComplementarioPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &complementario, derechoPecuniarioSolicitado)
							if errComplementarioPost != nil {
								logs.Error(errComplementarioPost)
								errorGetAll = true
							}
						}

					} else {
						errorGetAll = true
						logs.Error(errJson)
					}

				} else {
					errorGetAll = true
					logs.Error(errCodigo)
				}
			} else {
				errorGetAll = true
				logs.Error(errParam)
			}
		} else {
			errorGetAll = true
			logs.Error(err)
		}
	} else {
		errorGetAll = true
		logs.Error(err)
	}

	if !errorGetAll {
		return complementario, nil
	}
	return nil, errors.New("error del servicio PostGenerarDerechoPecuniarioEstudiante: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func GetEstadoRecibo(idPersona string, idPeriodo string) (interface{}, error) {
	var Recibos []map[string]interface{}
	var Periodo map[string]interface{}
	var ReciboXML map[string]interface{}
	var resultadoAux []map[string]interface{}
	resultado := make([]map[string]interface{}, 0)
	var Derecho map[string]interface{}
	var Programa map[string]interface{}
	var Solicitudes []map[string]interface{}
	var Estado string
	var PeriodoConsulta string
	var errorGetAll bool

	errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query=id:"+idPeriodo, &Periodo)
	if errPeriodo == nil {
		if Periodo != nil && fmt.Sprintf("%v", Periodo["Data"]) != "[map[]]" {
			PeriodoConsulta = fmt.Sprint(Periodo["Data"].([]interface{})[0].(map[string]interface{})["Year"])

			//Se consultan todos los recibos de derechos pecuniarios relacionados a ese tercero
			errRecibo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=0&query=InfoComplementariaId.Id:307,TerceroId.Id:"+idPersona, &Recibos)
			if errRecibo == nil {
				if Recibos != nil && fmt.Sprintf("%v", Recibos[0]) != "map[]" {
					// Ciclo for que recorre todos los recibos de derechos pecuniarios solicitados por el tercero
					resultadoAux = make([]map[string]interface{}, len(Recibos))
					for i := 0; i < len(Recibos); i++ {
						ReciboDerecho := "--"

						var reciboJson map[string]interface{}
						if err := json.Unmarshal([]byte(Recibos[i]["Dato"].(string)), &reciboJson); err == nil {
							ReciboDerecho = fmt.Sprintf("%v", reciboJson["Recibo"])
						}

						if strings.Split(ReciboDerecho, "/")[1] == PeriodoConsulta {
							errRecibo := request.GetJsonWSO2("http://"+beego.AppConfig.String("ConsultarReciboJbpmService")+"consulta_recibo/"+ReciboDerecho, &ReciboXML)
							if errRecibo == nil {
								if ReciboXML != nil && fmt.Sprintf("%v", ReciboXML) != "map[reciboCollection:map[]]" && fmt.Sprintf("%v", ReciboXML) != "map[]" {
									//Fecha límite de pago extraordinario
									Fecha := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_extraordinario"].(string)
									EstadoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["estado"]
									PagoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["pago"]
									Valor := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["valor_ordinario"]
									concepto := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["observaciones"]
									Fecha_pago := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_ordinario"]
									Cedula_estudiante := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["documento"]
									ProgramaAcademicoId := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["carrera"]
									IdConcepto := "0"

									switch concepto {
									case "CERTIFICADO DE NOTAS":
										IdConcepto = "40"
									case "DERECHOS DE GRADO":
										IdConcepto = "50"
									case "DUPLICADO DEL DIPLOMA DE GRADO":
										IdConcepto = "51"
									case "DUPLICADO DEL CARNET ESTUDIANTIL":
										IdConcepto = "44"
									case "CURSOS VACIONALES":
										IdConcepto = "31"
									case "CONSTANCIAS DE ESTUDIO":
										IdConcepto = "41"
									case "COPIA ACTA DE GRADO":
										IdConcepto = "49"
									case "CARNET ESTUDIANTIL":
										IdConcepto = "42"
									case "Inscripcion Virtual":
										conceptos := make([]string, 0)
										conceptos = append(conceptos, "40", "50", "51", "44", "31", "41", "49", "42")
										rand.Seed(time.Now().Unix())
										IdConcepto = conceptos[rand.Intn(len(conceptos))]
									}

									//Nombre del derecho pecuniario
									errDerecho := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId.CodigoAbreviacion:"+IdConcepto+",PeriodoId.Id:"+idPeriodo+",Activo:true", &Derecho)
									NombreConcepto := "---"
									if errDerecho == nil {
										if Derecho != nil && fmt.Sprintf("%v", Derecho["Data"]) != "map[]" {
											/////////////////////////////////////////////////////////////
											Resultado := Derecho["Data"].([]interface{})[0].(map[string]interface{})["Valor"].(string)

											var ResultadoJson map[string]interface{}
											if err := json.Unmarshal([]byte(Resultado), &ResultadoJson); err == nil {
												Valor = ResultadoJson["Costo"]
											}
											////////////////////////////////////////////////////////////
											NombreConcepto = fmt.Sprint(Derecho["Data"].([]interface{})[0].(map[string]interface{})["ParametroId"].(map[string]interface{})["Nombre"])
										} else {
											errorGetAll = true
										}

									} else {
										errorGetAll = true
									}

									valorPagado := ""
									fechaPago := ""
									var RespuestaDocID map[string]interface{}

									//Verificación si el recibo de pago se encuentra activo y pago
									if (EstadoRecibo == "A" && PagoRecibo == "S") || reciboJson["SolicitudId"] != "" {
										Estado = "Pago"

										valorPagado = fmt.Sprintf("%v", Valor)
										fechaPago = "" // Validar origen del dato

										//Información de la solicitud
										errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud?query=Id:"+fmt.Sprintf("%v", reciboJson["SolicitudId"]), &Solicitudes)
										if errSolicitud == nil {
											if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" && Solicitudes[0]["Resultado"] != nil {
												Resultado := Solicitudes[0]["Resultado"].(string)
												Estado = Solicitudes[0]["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Nombre"].(string)

												var ResultadoJson map[string]interface{}
												if err := json.Unmarshal([]byte(Resultado), &ResultadoJson); err == nil {
													RespuestaDocID = ResultadoJson
												}
											} else {
												errorGetAll = true
											}
										}

									} else {
										//Verifica si el recibo está vencido o no
										ATiempo, err := models.VerificarFechaLimite(Fecha)
										if err == nil {
											if ATiempo {
												Estado = "Pendiente pago"
											} else {
												Estado = "Vencido"
											}
										} else {
											Estado = "Vencido"
										}
									}

									errPrograma := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion/"+fmt.Sprintf("%v", ProgramaAcademicoId), &Programa)
									nombrePrograma := "---"
									if errPrograma == nil {
										nombrePrograma = fmt.Sprint(Programa["Nombre"])
									}

									resultadoAux[i] = map[string]interface{}{
										"Codigo":              IdConcepto,
										"Valor":               Valor,
										"Concepto":            NombreConcepto,
										"Id":                  strings.Split(ReciboDerecho, "/")[0],
										"FechaCreacion":       Recibos[i]["FechaCreacion"],
										"Estado":              Estado,
										"FechaOrdinaria":      Fecha_pago,
										"ProgramaAcademicoId": ProgramaAcademicoId,
										"ProgramaAcademico":   nombrePrograma,
										"Cedula_estudiante":   Cedula_estudiante,
										"Codigo_estudiante":   fmt.Sprintf("%v", reciboJson["CodigoAsociado"]),
										"Periodo":             PeriodoConsulta,
										"ValorPagado":         valorPagado,
										"FechaPago":           fechaPago,
										"VerRespuesta":        RespuestaDocID,
										"IdComplementario":    fmt.Sprintf("%v", Recibos[i]["Id"]),
									}

								} else {
									if len(resultado) > 0 {
										errorGetAll = false
									} else {
										errorGetAll = true
									}
								}
							} else {
								errorGetAll = true
							}
						}
						if fmt.Sprintf("%v", resultadoAux[i]) != "map[]" {
							resultado = append(resultado, resultadoAux[i])
							errorGetAll = false
						}
					}
				} else {
					errorGetAll = true
				}
			} else {
				errorGetAll = true
			}
		}
	}

	if !errorGetAll {
		return resultado, nil
	}
	return nil, errors.New("error del servicio GetEstadoRecibo: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func GetConsultarPersona(idPersona string) (interface{}, error) {
	var resultado map[string]interface{}
	var persona []map[string]interface{}
	var errorGetAll bool

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idPersona, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {

			var identificacion []map[string]interface{}

			errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId.Id:"+idPersona+"&sortby=Id&order=desc&limit=0", &identificacion)
			if errIdentificacion == nil && fmt.Sprintf("%v", identificacion[0]) != "map[]" {
				if identificacion[0]["Status"] != 404 {
					var codigos []map[string]interface{}
					var proyecto []map[string]interface{}

					resultado = persona[0]
					resultado["NumeroIdentificacion"] = identificacion[0]["Numero"]
					resultado["TipoIdentificacion"] = identificacion[0]["TipoDocumentoId"]
					resultado["FechaExpedicion"] = identificacion[0]["FechaExpedicion"]
					resultado["SoporteDocumento"] = identificacion[0]["DocumentoSoporte"]

					errCodigoEst := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.Id:93&limit=0", &codigos)
					if errCodigoEst == nil && fmt.Sprintf("%v", codigos[0]) != "map[]" {
						for _, codigo := range codigos {
							errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Codigo:"+codigo["Dato"].(string)[5:8], &proyecto)
							if errProyecto == nil && fmt.Sprintf("%v", proyecto[0]) != "map[]" {
								codigo["Proyecto"] = codigo["Dato"].(string) + " Proyecto: " + codigo["Dato"].(string)[5:8] + " - " + proyecto[0]["Nombre"].(string)
								codigo["IdProyecto"] = proyecto[0]["Codigo"]
							}
						}

						resultado["Codigos"] = codigos
					}

				} else {
					if identificacion[0]["Message"] == "Not found resource" {
						errorGetAll = true
					} else {
						logs.Error(identificacion)
						errorGetAll = true
					}
				}
			} else {
				logs.Error(identificacion)
				errorGetAll = true
			}
		} else {
			if persona[0]["Message"] == "Not found resource" {
				errorGetAll = true
			} else {
				logs.Error(persona)
				errorGetAll = true
			}
		}
	} else {
		logs.Error(persona)
		errorGetAll = true
	}

	if !errorGetAll {
		return resultado, nil
	}
	return nil, errors.New("error del servicio GetConsultarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func PostSolicitudDerechoPecuniario(data []byte) (interface{}, error) {
	var Referencia string
	var resDocs []interface{}
	var SolicitudPost map[string]interface{}
	var SolicitantePost map[string]interface{}
	var SolicitudEvolucionEstadoPost map[string]interface{}
	var SolicitudData map[string]interface{}
	var Infocomplementario map[string]interface{}
	var InfocomplementarioPut map[string]interface{}
	var Derecho map[string]interface{}
	var TerceroSolicitante map[string]interface{}
	resultado := make(map[string]interface{})
	var errorGetAll bool

	if err := json.Unmarshal(data, &SolicitudData); err == nil {

		auxDoc := []map[string]interface{}{}
		documento := map[string]interface{}{
			"IdTipoDocumento": SolicitudData["comprobanteRecibo"].(map[string]interface{})["IdTipoDocumento"],
			"nombre":          SolicitudData["comprobanteRecibo"].(map[string]interface{})["nombre"],
			"metadatos":       SolicitudData["comprobanteRecibo"].(map[string]interface{})["metadatos"],
			"descripcion":     SolicitudData["comprobanteRecibo"].(map[string]interface{})["descripcion"],
			"file":            SolicitudData["comprobanteRecibo"].(map[string]interface{})["file"],
		}
		auxDoc = append(auxDoc, documento)

		doc, errDoc := models.RegistrarDoc(auxDoc)
		if errDoc == nil {
			docTem := map[string]interface{}{
				"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
				"Enlace":        doc.(map[string]interface{})["Enlace"],
				"Id":            doc.(map[string]interface{})["Id"],
				"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
				"Activo":        doc.(map[string]interface{})["Activo"],
			}

			resDocs = append(resDocs, docTem)
		}

		var jsonTerceroSolicitante []byte
		var jsonDerechoPecuniarioId []byte

		errTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", SolicitudData["SolicitanteId"]), &TerceroSolicitante)
		if errTercero == nil && TerceroSolicitante != nil {
			if fmt.Sprintf("%v", TerceroSolicitante) != "map[]" && TerceroSolicitante["Status"] != "404" {
				jsonTerceroSolicitante, _ = json.Marshal(TerceroSolicitante)
			}
		}

		errDerecho := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId.CodigoAbreviacion:"+fmt.Sprintf("%v", SolicitudData["Codigo"])+",PeriodoId.Year:"+fmt.Sprintf("%v", SolicitudData["Periodo"])+",Activo:true", &Derecho)
		if errDerecho == nil && fmt.Sprintf("%v", Derecho["Data"]) != "map[]" {
			jsonDerechoPecuniarioId, _ = json.Marshal(Derecho["Data"].([]interface{})[0])
		}

		jsonDocumento, _ := json.Marshal(resDocs)
		jsonCodigoEstudiante, _ := json.Marshal(SolicitudData["Codigo_estudiante"])

		Referencia = "{\n\"DocSoportePago\": " + fmt.Sprintf("%v", string(jsonDocumento)) +
			",\n\"TerceroSolicitante\": " + fmt.Sprintf("%v", string(jsonTerceroSolicitante)) +
			",\n\"CodigoEstudiante\": " + fmt.Sprintf("%v", string(jsonCodigoEstudiante)) +
			",\n\"DerechoPecuniarioId\": " + fmt.Sprintf("%v", string(jsonDerechoPecuniarioId)) + "\n}"

		IdEstadoTipoSolicitud := 41

		SolicitudPracticas := map[string]interface{}{
			"EstadoTipoSolicitudId": map[string]interface{}{"Id": IdEstadoTipoSolicitud},
			"Referencia":            Referencia,
			"Resultado":             "",
			"FechaRadicacion":       fmt.Sprintf("%v", SolicitudData["FechaCreacion"]),
			"Activo":                true,
			"SolicitudPadreId":      nil,
		}

		errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud", "POST", &SolicitudPost, SolicitudPracticas)
		if errSolicitud == nil {
			if SolicitudPost["Success"] != false && fmt.Sprintf("%v", SolicitudPost) != "map[]" {
				resultado["Solicitud"] = SolicitudPost["Data"]
				IdSolicitud := SolicitudPost["Data"].(map[string]interface{})["Id"]

				//POST tabla solicitante
				Solicitante := map[string]interface{}{
					"TerceroId": SolicitudData["SolicitanteId"],
					"SolicitudId": map[string]interface{}{
						"Id": IdSolicitud,
					},
					"Activo": true,
				}

				errSolicitante := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante", "POST", &SolicitantePost, Solicitante)
				if errSolicitante == nil && fmt.Sprintf("%v", SolicitantePost["Status"]) != "400" {
					if SolicitantePost != nil && fmt.Sprintf("%v", SolicitantePost) != "map[]" {
						//POST a la tabla solicitud_evolucion estado
						SolicitudEvolucionEstado := map[string]interface{}{
							"TerceroId": SolicitudData["SolicitanteId"],
							"SolicitudId": map[string]interface{}{
								"Id": IdSolicitud,
							},
							"EstadoTipoSolicitudIdAnterior": nil,
							"EstadoTipoSolicitudId": map[string]interface{}{
								"Id": IdEstadoTipoSolicitud,
							},
							"Activo":      true,
							"FechaLimite": fmt.Sprintf("%v", SolicitudData["FechaCreacion"]),
						}

						errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
						if errSolicitudEvolucionEstado == nil {
							if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {
								idComplementario := SolicitudData["IdComplementario"]
								errInfoComplementario := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%v", idComplementario), &Infocomplementario)
								if errInfoComplementario == nil && Infocomplementario != nil {
									if fmt.Sprintf("%v", Infocomplementario) != "map[]" && Infocomplementario["Status"] != "404" {

										var InfocomplementarioJson map[string]interface{}
										if err := json.Unmarshal([]byte(Infocomplementario["Dato"].(string)), &InfocomplementarioJson); err == nil {
											Infocomplementario["Dato"] = `{"Recibo":` + `"` + fmt.Sprintf("%v", InfocomplementarioJson["Recibo"]) + `", ` + `"CodigoAsociado": "` + fmt.Sprintf("%v", InfocomplementarioJson["CodigoAsociado"]) + `", "SolicitudId":"` + fmt.Sprintf("%v", IdSolicitud) + `"}`

											errActuComplementario := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%v", idComplementario), "PUT", &InfocomplementarioPut, Infocomplementario)
											if errActuComplementario != nil {
												resultado["Solicitante"] = SolicitantePost["Data"]
											}
										}
									}
								}

							} else {
								errorGetAll = true
							}
						} else {
							var resultado2 map[string]interface{}
							request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
							request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante/"+fmt.Sprintf("%v", SolicitantePost["Id"]), "DELETE", &resultado2, nil)
							errorGetAll = true
						}
					} else {
						errorGetAll = true
					}
				} else {
					//Se elimina el registro de solicitud si no se puede hacer el POST a la tabla solicitante
					var resultado2 map[string]interface{}
					request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
					errorGetAll = true
				}
			} else {
				errorGetAll = true
			}
		} else {
			errorGetAll = true
		}
	} else {
		errorGetAll = true
	}

	if !errorGetAll {
		return resultado, nil
	}
	return nil, errors.New("error del servicio PostSolicitudDerechoPecuniario: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func GetSolicitudDerechoPecuniario() (interface{}, error) {
	var Solicitudes []map[string]interface{}
	var DatosIdentificacion []map[string]interface{}
	resultado := make([]map[string]interface{}, 0)
	var errorGetAll bool

	errSolicitudes := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud?query=EstadoTipoSolicitudId.Id:41,Activo:true&limit=0", &Solicitudes)
	if errSolicitudes == nil {
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" && Solicitudes[0]["Resultado"] != nil {
			for _, solicitud := range Solicitudes {
				referencia := solicitud["Referencia"].(string)
				FechaCreacion := fmt.Sprintf("%v", solicitud["FechaCreacion"])

				var referenciaJson map[string]interface{}
				if err := json.Unmarshal([]byte(referencia), &referenciaJson); err == nil {
					VerSoporte := fmt.Sprintf("%v", referenciaJson["DocSoportePago"].([]interface{})[0].(map[string]interface{})["Id"])
					TerceroSolicitanteId := fmt.Sprintf("%v", referenciaJson["TerceroSolicitante"].(map[string]interface{})["Id"])
					Nombre := fmt.Sprintf("%v", referenciaJson["TerceroSolicitante"].(map[string]interface{})["NombreCompleto"])
					Codigo := fmt.Sprintf("%v", referenciaJson["CodigoEstudiante"])
					DerechoPecuniarioId := referenciaJson["DerechoPecuniarioId"]
					DerechoValor := referenciaJson["DerechoPecuniarioId"].(map[string]interface{})["Valor"].(string)
					var valorJson map[string]interface{}
					valor := "0"
					if err := json.Unmarshal([]byte(DerechoValor), &valorJson); err == nil {
						valor = fmt.Sprintf("%v", valorJson["Costo"])
					}

					errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?limit=0&query=TerceroId.Id:"+TerceroSolicitanteId+",Activo:True", &DatosIdentificacion)
					if errIdentificacion == nil {
						if DatosIdentificacion != nil && fmt.Sprintf("%v", DatosIdentificacion[0]) != "map[]" {
							NombreIdentificacion := fmt.Sprintf("%v", DatosIdentificacion[0]["TipoDocumentoId"].(map[string]interface{})["Nombre"])
							CAIdentificacion := fmt.Sprintf("%v", DatosIdentificacion[0]["TipoDocumentoId"].(map[string]interface{})["CodigoAbreviacion"])
							Identificacion := NombreIdentificacion + " - " + CAIdentificacion

							resultadoAux := map[string]interface{}{
								"FechaCreacion":        FechaCreacion,
								"VerSoporte":           VerSoporte,
								"Nombre":               Nombre,
								"Codigo":               Codigo,
								"DerechoPecuniarioId":  DerechoPecuniarioId,
								"NombreIdentificacion": NombreIdentificacion,
								"CAIdentificacion":     CAIdentificacion,
								"Identificacion":       Identificacion,
								"Concepto":             fmt.Sprintf("%v", DerechoPecuniarioId.(map[string]interface{})["ParametroId"].(map[string]interface{})["Nombre"]),
								"Valor":                valor,
								"Id":                   fmt.Sprintf("%v", solicitud["Id"]),
								"Estado":               fmt.Sprintf("%v", solicitud["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Nombre"]),
							}

							resultado = append(resultado, resultadoAux)
						} else {
							errorGetAll = true
						}
					} else {
						errorGetAll = true
					}
				} else {
					errorGetAll = true
				}
			}
		} else {
			errorGetAll = true
		}
	} else {
		errorGetAll = true
	}

	if !errorGetAll {
		return resultado, nil
	}
	return nil, errors.New("error del servicio GetSolicitudDerechoPecuniario: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func PostRespuestaSolicitudDerechoPecuniario(idSolicitud string, data []byte) (interface{}, error) {
	var RespuestaSolicitud map[string]interface{}
	var Solicitud map[string]interface{}
	var SolicitudPut map[string]interface{}
	var anteriorEstado []map[string]interface{}
	var Resultado string

	var SolicitudEvolucionEstadoPost map[string]interface{}
	var anteriorEstadoPost map[string]interface{}
	var resDocs []interface{}
	var errorGetAll bool

	if err := json.Unmarshal(data, &RespuestaSolicitud); err == nil {

		// Consulta de información de la solicitud
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+idSolicitud, &Solicitud)
		if errSolicitud == nil {
			if Solicitud != nil && fmt.Sprintf("%v", Solicitud["Status"]) != "404" {

				if RespuestaSolicitud["DocRespuesta"] != nil {
					if len(RespuestaSolicitud["DocRespuesta"].([]interface{})) > 0 {
						for i := range RespuestaSolicitud["DocRespuesta"].([]interface{}) {
							auxDoc := []map[string]interface{}{}
							documento := map[string]interface{}{
								"IdTipoDocumento": RespuestaSolicitud["DocRespuesta"].([]interface{})[i].(map[string]interface{})["IdTipoDocumento"],
								"nombre":          RespuestaSolicitud["DocRespuesta"].([]interface{})[i].(map[string]interface{})["nombre"],
								"metadatos":       RespuestaSolicitud["DocRespuesta"].([]interface{})[i].(map[string]interface{})["metadatos"],
								"descripcion":     RespuestaSolicitud["DocRespuesta"].([]interface{})[i].(map[string]interface{})["descripcion"],
								"file":            RespuestaSolicitud["DocRespuesta"].([]interface{})[i].(map[string]interface{})["file"],
							}
							auxDoc = append(auxDoc, documento)

							doc, errDoc := models.RegistrarDoc(auxDoc)
							if errDoc == nil {
								docTem := map[string]interface{}{
									"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
									"Enlace":        doc.(map[string]interface{})["Enlace"],
									"Id":            doc.(map[string]interface{})["Id"],
									"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
									"Activo":        doc.(map[string]interface{})["Activo"],
								}

								resDocs = append(resDocs, docTem)
							}
						}
					}

					jsonDocumento, _ := json.Marshal(resDocs)
					jsonTerceroResponsable, _ := json.Marshal(RespuestaSolicitud["TerceroResponasble"])

					Resultado = "{\n\"DocRespuesta\": " + fmt.Sprintf("%v", string(jsonDocumento)) +
						",\n\"Observacion\": \"" + fmt.Sprintf("%v", RespuestaSolicitud["Observacion"]) + "\"" +
						",\n\"TerceroResponasble\": " + fmt.Sprintf("%v", string(jsonTerceroResponsable)) + "\n}"
				}

				// Actualización del anterior estado
				errAntEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=activo:true,solicitudId.Id:"+idSolicitud, &anteriorEstado)
				if errAntEstado == nil {
					if anteriorEstado != nil && fmt.Sprintf("%v", anteriorEstado) != "map[]" {

						anteriorEstado[0]["Activo"] = false
						estadoAnteriorId := fmt.Sprintf("%v", anteriorEstado[0]["Id"])

						errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado/"+estadoAnteriorId, "PUT", &anteriorEstadoPost, anteriorEstado[0])
						if errSolicitudEvolucionEstado == nil {

							id, _ := strconv.Atoi(idSolicitud)
							SolicitudEvolucionEstado := map[string]interface{}{
								"TerceroId": int(RespuestaSolicitud["TerceroResponasble"].(map[string]interface{})["Id"].(float64)),
								"SolicitudId": map[string]interface{}{
									"Id": id,
								},
								"EstadoTipoSolicitudId": map[string]interface{}{
									"Id": 42,
								},
								"EstadoTipoSolicitudIdAnterior": map[string]interface{}{
									"Id": 41,
								},
								"Activo":      true,
								"FechaLimite": RespuestaSolicitud["FechaRespuesta"],
							}

							errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
							if errSolicitudEvolucionEstado == nil {
								if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {

									Solicitud["Resultado"] = Resultado
									Solicitud["EstadoTipoSolicitudId"] = SolicitudEvolucionEstadoPost["Data"].(map[string]interface{})["EstadoTipoSolicitudId"]
									// Solicitud["EstadoTipoSolicitudId"].(map[string]interface{})["Activo"] = true
									Solicitud["SolicitudFinalizada"] = true

									errPutEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+idSolicitud, "PUT", &SolicitudPut, Solicitud)

									if errPutEstado == nil {
										if SolicitudPut["Status"] == "400" {
											errorGetAll = true
										}
									} else {
										errorGetAll = true
									}

								} else {
									errorGetAll = true
								}
							} else {
								errorGetAll = true
							}

						} else {
							errorGetAll = true
						}

					} else {
						errorGetAll = true
					}

				} else {
					errorGetAll = true
				}

			} else {
				errorGetAll = true
			}

		} else {
			errorGetAll = true
		}

	} else {
		errorGetAll = true
	}

	if !errorGetAll {
		return SolicitudPut, nil
	}
	return nil, errors.New("error del servicio PostRespuestaSolicitudDerechoPecuniario: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}
