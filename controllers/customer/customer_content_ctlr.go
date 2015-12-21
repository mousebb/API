package customerCtlr

//
// // Get it all
// func GetAllContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	content, err := custcontent.AllCustomerContent(dtx.APIKey)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting customer content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// // Get Content by Content Id
// // Returns: CustomerContent
// func GetContentById(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting customer content ID", err, rw, r)
// 		return ""
// 	}
//
// 	content, err := custcontent.GetCustomerContent(id, dtx.APIKey)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting customer content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// // Get Content by Content Id
// // Returns: CustomerContent
// func GetContentRevisionsById(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting customer content ID", err, rw, r)
// 		return ""
// 	}
//
// 	revs, err := custcontent.GetCustomerContentRevisions(id, dtx.APIKey)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting customer content revisions", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(revs))
// }
//
// // Part Content Endpoints
// func AllPartContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	content, err := custcontent.GetAllPartContent(dtx.APIKey)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting all part content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// func UniquePartContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	partID, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting part ID", err, rw, r)
// 		return ""
// 	}
//
// 	content, err := custcontent.GetPartContent(partID, dtx.APIKey)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting part content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// func CreatePartContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting part ID", err, rw, r)
// 		return ""
// 	}
//
// 	// Defer the body closing until we're finished
// 	defer r.Body.Close()
//
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		apierror.GenerateError("Trouble reading request body while creating part content", err, rw, r)
// 		return ""
// 	}
//
// 	var content custcontent.CustomerContent
// 	if err = json.Unmarshal(body, &content); err != nil {
// 		apierror.GenerateError("Trouble unmarshalling json request body while creating part content", err, rw, r)
// 		return ""
// 	}
//
// 	if err = content.Save(id, 0, dtx.APIKey); err != nil {
// 		apierror.GenerateError("Trouble creating part content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
// func UpdatePartContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting part ID", err, rw, r)
// 		return ""
// 	}
//
// 	// Defer the body closing until we're finished
// 	defer r.Body.Close()
//
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		apierror.GenerateError("Trouble reading request body while updating part content", err, rw, r)
// 		return ""
// 	}
//
// 	var content custcontent.CustomerContent
// 	if err = json.Unmarshal(body, &content); err != nil {
// 		apierror.GenerateError("Trouble unmarshalling json request body while updating part content", err, rw, r)
// 		return ""
// 	}
//
// 	if err = content.Save(id, 0, dtx.APIKey); err != nil {
// 		apierror.GenerateError("Trouble updating part content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// func DeletePartContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting part ID", err, rw, r)
// 		return ""
// 	}
//
// 	// Defer the body closing until we're finished
// 	defer r.Body.Close()
//
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		apierror.GenerateError("Trouble reading request body while deleting part content", err, rw, r)
// 		return ""
// 	}
//
// 	var content custcontent.CustomerContent
// 	if err = json.Unmarshal(body, &content); err != nil {
// 		apierror.GenerateError("Trouble unmarshalling json request body while deleting part content", err, rw, r)
// 		return ""
// 	}
//
// 	if err = content.Delete(id, 0, dtx.APIKey); err != nil {
// 		apierror.GenerateError("Trouble deleting part content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// // Category Content Endpoints
// func AllCategoryContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	content, err := custcontent.GetAllCategoryContent(dtx.APIKey)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting all category content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// func UniqueCategoryContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	catID, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting category ID", err, rw, r)
// 		return ""
// 	}
//
// 	content, err := custcontent.GetCategoryContent(catID, dtx.APIKey)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting category content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// func CreateCategoryContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting category ID", err, rw, r)
// 		return ""
// 	}
//
// 	// Defer the body closing until we're finished
// 	defer r.Body.Close()
//
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		apierror.GenerateError("Trouble reading request body while creating category content", err, rw, r)
// 		return ""
// 	}
//
// 	var content custcontent.CustomerContent
// 	if err = json.Unmarshal(body, &content); err != nil {
// 		apierror.GenerateError("Trouble unmarshalling json request body while creating category content", err, rw, r)
// 		return ""
// 	}
//
// 	if err = content.Save(0, id, dtx.APIKey); err != nil {
// 		apierror.GenerateError("Trouble creating category content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// func UpdateCategoryContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting category ID", err, rw, r)
// 		return ""
// 	}
//
// 	// Defer the body closing until we're finished
// 	defer r.Body.Close()
//
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		apierror.GenerateError("Trouble reading request body while updating category content", err, rw, r)
// 		return ""
// 	}
//
// 	var content custcontent.CustomerContent
// 	if err = json.Unmarshal(body, &content); err != nil {
// 		apierror.GenerateError("Trouble unmarshalling json request body while updating category content", err, rw, r)
// 		return ""
// 	}
//
// 	if err = content.Save(0, id, dtx.APIKey); err != nil {
// 		apierror.GenerateError("Trouble updating category content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// func DeleteCategoryContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting category ID", err, rw, r)
// 		return ""
// 	}
//
// 	// Defer the body closing until we're finished
// 	defer r.Body.Close()
//
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		apierror.GenerateError("Trouble reading request body while deleting category content", err, rw, r)
// 		return ""
// 	}
//
// 	var content custcontent.CustomerContent
// 	if err = json.Unmarshal(body, &content); err != nil {
// 		apierror.GenerateError("Trouble unmarshalling json request body while deleting category content", err, rw, r)
// 		return ""
// 	}
//
// 	if err = content.Delete(0, id, dtx.APIKey); err != nil {
// 		apierror.GenerateError("Trouble deleting category content", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(content))
// }
//
// // Content Types
// func GetAllContentTypes(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	types, err := custcontent.AllCustomerContentTypes()
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting all content types", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(types))
// }
