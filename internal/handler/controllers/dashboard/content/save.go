package edit

func SaveContentHandler() {
	//	cid := req.FormValue("id")
	//	t := req.FormValue("type")

	//	contentType, ok := contentTypes[t]
	//	if !ok {
	//		_, err := fmt.Fprintf(res, contentPkg.ErrTypeNotRegistered.Error(), t)
	//		if err != nil {
	//			log.WithField("Error", err).Warning("Failed to write response")
	//		}

	//		return
	//	}

	//	err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
	//	if err != nil {
	//		log.WithField("Error", err).Warning("Failed to parse form")
	//		r.Renderer().InternalServerError(res)
	//		return
	//	}

	//	files, err := request.GetRequestFiles(req)
	//	if err != nil {
	//		log.WithField("Error", err).Warning("Failed to get request files")
	//		r.Renderer().InternalServerError(res)
	//		return
	//	}

	//	urlPaths, err := storageService.StoreFiles(files)
	//	if err != nil {
	//		log.WithField("Error", err).Warning("Failed to get save files")
	//		return
	//	}

	//	for name, urlPath := range urlPaths {
	//		req.PostForm.Set(name, urlPath)
	//	}

	//	entity, err := request.GetEntityFromFormData(contentType, req.PostForm)
	//	if err != nil {
	//		log.WithField("Error", err).Warning("Failed to map request entity")
	//		return
	//	}

	//	hook, ok := entity.(item.Hookable)
	//	if !ok {
	//		log.Println("Type", t, "does not implement item.Hookable or embed item.Item.")
	//		r.Renderer().BadRequest(res)
	//		return
	//	}

	//	if cid == "" {
	//		err = hook.BeforeAdminCreate(res, req)
	//		if err != nil {
	//			log.Println("Error running BeforeAdminCreate method in editHandler for:", t, err)
	//			return
	//		}
	//	} else {
	//		err = hook.BeforeAdminUpdate(res, req)
	//		if err != nil {
	//			log.Println("Error running BeforeAdminUpdate method in editHandler for:", t, err)
	//			return
	//		}
	//	}

	//	err = hook.BeforeSave(res, req)
	//	if err != nil {
	//		log.Println("Error running BeforeSave method in editHandler for:", t, err)
	//		return
	//	}

	//	if cid == "" {
	//		cid, err = contentService.CreateContent(t, entity)
	//		if err != nil {
	//			log.WithField("Error", err).Warning("Failed to create content")
	//			r.Renderer().InternalServerError(res)
	//			return
	//		}

	//	} else {
	//		_, err = contentService.UpdateContent(t, cid, entity)
	//		if err != nil {
	//			log.WithField("Error", err).Warning("Failed to update content")
	//			r.Renderer().InternalServerError(res)
	//			return
	//		}
	//	}

	//	// set the target in the context so user can get saved value from db in hook
	//	ctx := context.WithValue(req.Context(), "target", fmt.Sprintf("%s:%s", t, cid))
	//	req = req.WithContext(ctx)

	//	err = hook.AfterSave(res, req)
	//	if err != nil {
	//		log.Println("Error running AfterSave method in editHandler for:", t, err)
	//		return
	//	}

	//	if cid == "" {
	//		err = hook.AfterAdminCreate(res, req)
	//		if err != nil {
	//			log.Println("Error running AfterAdminUpdate method in editHandler for:", t, err)
	//			return
	//		}
	//	} else {
	//		err = hook.AfterAdminUpdate(res, req)
	//		if err != nil {
	//			log.Println("Error running AfterAdminUpdate method in editHandler for:", t, err)
	//			return
	//		}
	//	}

	//	r.Redirect(req, res, req.URL.Path+"?type="+t+"&id="+cid)
}
