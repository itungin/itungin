package route

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/controller"
	"github.com/gocroot/helper/at"
	// "github.com/gocroot/helper/at"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}
	config.SetEnv()

	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		controller.GetHome(w, r)
	//gis
	case method == "POST" && path == "/data/gis/lokasi":
		controller.GetRegion(w, r)
	//chat bot inbox
	case method == "POST" && at.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)
	//masking list nmor official
	case method == "GET" && path == "/data/phone/all":
		controller.GetBotList(w, r)
	//akses data helpdesk layanan user
	case method == "GET" && path == "/data/user/helpdesk/all":
		controller.GetHelpdeskAll(w, r)
	case method == "GET" && path == "/data/user/helpdesk/masuk":
		controller.GetLatestHelpdeskMasuk(w, r)
	case method == "GET" && path == "/data/user/helpdesk/selesai":
		controller.GetLatestHelpdeskSelesai(w, r)
	//pamong desa data from api
	case method == "GET" && path == "/data/lms/user":
		controller.GetDataUserFromApi(w, r)
	//simpan testimoni dari pamong desa lms api
	case method == "POST" && path == "/data/lms/testi":
		controller.PostTestimoni(w, r)
		//get random 4 testi
	case method == "GET" && path == "/data/lms/random/testi":
		controller.GetRandomTesti4(w, r)
	//mendapatkan data sent item
	case method == "GET" && at.URLParam(path, "/data/peserta/sent/:id"):
		controller.GetSentItem(w, r)
	//simpan feedback unsubs user
	case method == "POST" && path == "/data/peserta/unsubscribe":
		controller.PostUnsubscribe(w, r)
	//generate token linked device
	case method == "PUT" && path == "/data/user":
		controller.PutTokenDataUser(w, r)
	//Menambhahkan data nomor sender untuk broadcast
	case method == "PUT" && path == "/data/sender":
		controller.PutNomorBlast(w, r)
	//mendapatkan data list nomor sender untuk broadcast
	case method == "GET" && path == "/data/sender":
		controller.GetDataSenders(w, r)
	//mendapatkan data list nomor sender yang kena blokir dari broadcast
	case method == "GET" && path == "/data/blokir":
		controller.GetDataSendersTerblokir(w, r)
	//mendapatkan data rekap pengiriman wa blast
	case method == "GET" && path == "/data/rekap":
		controller.GetRekapBlast(w, r)
	//mendapatkan data faq
	case method == "GET" && at.URLParam(path, "/data/faq/:id"):
		controller.GetFAQ(w, r)
	//legacy
	case method == "PUT" && path == "/data/user/task/doing":
		controller.PutTaskUser(w, r)
	case method == "GET" && path == "/data/user/task/done":
		controller.GetTaskDone(w, r)
	case method == "POST" && path == "/data/user/task/done":
		controller.PostTaskUser(w, r)
	case method == "GET" && path == "/data/pushrepo/kemarin":
		controller.GetYesterdayDistincWAGroup(w, r)

	//helpdesk
	//mendapatkan data tiket
	case method == "GET" && at.URLParam(path, "/data/tiket/closed/:id"):
		controller.GetClosedTicket(w, r)
	//simpan feedback tiket user
	case method == "POST" && path == "/data/tiket/rate":
		controller.PostMasukanTiket(w, r)
		// order
	case method == "POST" && at.URLParam(path, "/data/order/:namalapak"):
		controller.HandleOrder(w, r)
	//user data
	case method == "GET" && path == "/data/user":
		controller.GetDataUser(w, r)
	//user pendaftaran
	case method == "POST" && path == "/auth/register/users": //mendapatkan email gmail
		controller.RegisterGmailAuth(w, r)
	case method == "POST" && path == "/data/user":
		controller.PostDataUser(w, r)
	case method == "POST" && path == "/upload/profpic": //upload gambar profile
		controller.UploadProfilePictureHandler(w, r)
	case method == "POST" && path == "/data/user/bio":
		controller.PostDataBioUser(w, r)
		/* 	case method == "POST" && at.URLParam(path, "/data/user/wa/:nomorwa"):
		controller.PostDataUserFromWA(w, r) */
	//data proyek
	case method == "GET" && path == "/data/proyek":
		controller.GetDataProject(w, r)
	case method == "GET" && path == "/data/proyek/approved": //akses untuk manager
		controller.GetEditorApprovedProject(w, r)
	case method == "POST" && path == "/data/proyek":
		controller.PostDataProject(w, r)
	case method == "PUT" && path == "/data/metadatabuku":
		controller.PutMetaDataProject(w, r)
	case method == "PUT" && path == "/data/proyek/publishbuku": //publish buku isbn by manager
		controller.PutPublishProject(w, r)
	case method == "PUT" && path == "/data/proyek":
		controller.PutDataProject(w, r)
	case method == "DELETE" && path == "/data/proyek":
		controller.DeleteDataProject(w, r)
	case method == "GET" && path == "/data/proyek/anggota":
		controller.GetDataMemberProject(w, r)
	case method == "GET" && path == "/data/proyek/editor":
		controller.GetDataEditorProject(w, r)
	case method == "DELETE" && path == "/data/proyek/anggota":
		controller.DeleteDataMemberProject(w, r)
	case method == "POST" && path == "/data/proyek/anggota":
		controller.PostDataMemberProject(w, r)
	case method == "POST" && path == "/data/proyek/editor": //set editor oleh owner
		controller.PostDataEditorProject(w, r)
	case method == "PUT" && path == "/data/proyek/editor": //set approved oleh editor
		controller.PUtApprovedEditorProject(w, r)
	//upload cover,draft,pdf,sampul buku project
	case method == "POST" && at.URLParam(path, "/upload/coverbuku/:projectid"):
		controller.UploadCoverBukuWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/draftbuku/:projectid"):
		controller.UploadDraftBukuWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/draftpdfbuku/:projectid"):
		controller.UploadDraftBukuPDFWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/sampulpdfbuku/:projectid"):
		controller.UploadSampulBukuPDFWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/spk/:projectid"):
		controller.UploadSPKPDFWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/spi/:projectid"):
		controller.UploadSPIPDFWithParamFileHandler(w, r)
	case method == "GET" && at.URLParam(path, "/download/draft/:path"): //downoad file draft
		controller.AksesFileRepoDraft(w, r)
	case method == "POST" && path == "/data/proyek/katalog": //post blog katalog
		controller.PostKatalogBuku(w, r)
	case method == "GET" && at.URLParam(path, "/download/dokped/spk/:namaproject"): //base64 namaproject
		controller.GetFileDraftSPK(w, r)
	case method == "GET" && at.URLParam(path, "/download/dokped/spkt/:namaproject"): //base64 namaproject
		controller.GetFileDraftSPKT(w, r)
	case method == "GET" && at.URLParam(path, "/download/dokped/spi/:path"): //base64 path sampul
		controller.GetFileDraftSPI(w, r)

	case method == "POST" && path == "/data/proyek/menu":
		controller.PostDataMenuProject(w, r)
	case method == "POST" && path == "/approvebimbingan":
		controller.ApproveBimbinganbyPoin(w, r)
	case method == "DELETE" && path == "/data/proyek/menu":
		controller.DeleteDataMenuProject(w, r)
	case method == "POST" && path == "/notif/ux/postlaporan":
		controller.PostLaporan(w, r)
	case method == "POST" && path == "/notif/ux/postfeedback":
		controller.PostFeedback(w, r)

	case method == "POST" && path == "/notif/ux/postmeeting":
		controller.PostMeeting(w, r)
	case method == "POST" && at.URLParam(path, "/notif/ux/postpresensi/:id"):
		controller.PostPresensi(w, r)
	case method == "POST" && at.URLParam(path, "/notif/ux/posttasklists/:id"):
		controller.PostTaskList(w, r)
	case method == "POST" && at.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)
	// LMS
	case method == "GET" && path == "/lms/refresh/cookie":
		controller.RefreshLMSCookie(w, r)
	case method == "GET" && path == "/lms/count/user":
		controller.GetCountDocUser(w, r)
	// Google Auth
	case method == "POST" && path == "/auth/users":
		controller.Auth(w, r)
	case method == "POST" && path == "/auth/login":
		controller.GeneratePasswordHandler(w, r)
	case method == "POST" && path == "/auth/verify":
		controller.VerifyPasswordHandler(w, r)
	case method == "POST" && path == "/auth/resend":
		controller.ResendPasswordHandler(w, r)
		// Produk
	case method == "POST" && path == "/products":
		controller.CreateProduct(w, r)
	case method == "GET" && path == "/products":
		controller.GetProducts(w, r)
	case method == "GET" && path == "/product-id":
		controller.GetProductByID(w, r)
	case method == "PUT" && path == "/products":
		controller.UpdateProduct(w, r)
	case method == "DELETE" && path == "/products/{id}":
		controller.DeleteProduct(w, r)
	case method == "GET" && path == "/products-export-csv":
		controller.ExportProductsToCSV(w, r)
				// Expense
	case method == "POST" && path == "/expense":
		controller.CreateExpenseTransaction(w, r)
	case method == "GET" && path == "/expense":
		controller.GetExpenses(w, r)
	case method == "GET" && path == "/expense/{id}":
		controller.GetExpenseByID(w, r)
	case method == "PUT" && path == "/expense/{id}":
		controller.UpdateExpense(w, r)
	case method == "DELETE" && path == "/expense/{id}":
		controller.DeleteExpense(w, r)
	case method == "GET" && path == "/expense-export-csv":
		controller.ExportProductsToCSV(w, r)
	// Sales
	case method == "POST" && path == "/sales":
		controller.CreateSalesTransaction(w, r)
	case method == "GET" && path == "/sales":
		controller.GetSalesTransactions(w, r)
	case method == "GET" && path == "/sales/{id}":
		controller.GetSalesTransactionByID(w, r)
	case method == "PUT" && path == "/sales/{id}":
		controller.UpdateSalesTransaction(w, r)
	case method == "DELETE" && path == "/sales/{id}":
		controller.DeleteSalesTransaction(w, r)
	case method == "GET" && path == "/sales-export-csv":
		controller.ExportProductsToCSV(w, r)
	// Pelanggan
	case method == "POST" && path == "/pelanggan":
		controller.CreateCustomer(w, r)
	case method == "GET" && path == "/pelanggan":
		controller.GetCustomers(w, r)
	case method == "GET" && path == "/pelanggan/{id}":
		controller.GetCustomerByID(w, r)
	case method == "PUT" && path == "/pelanggan/{id}":
		controller.UpdateCustomer(w, r)
	case method == "DELETE" && path == "/pelanggan/{id}":
		controller.DeleteCustomer(w, r)
	// Laporan Akuntan
	case method == "POST" && path == "/laporan":
		controller.CreateFinancialReport(w, r)
	case method == "GET" && path == "/laporan":
		controller.GetFinancialReports(w, r)
	case method == "GET" && path == "/laporan/{id}":
		controller.GetFinancialReportByID(w, r)
	case method == "DELETE" && path == "/laporan/{id}":
		controller.DeleteFinancialReport(w, r)

	// Google Auth
	default:
		controller.NotFound(w, r)
	}
}


// func URL(w http.ResponseWriter, r *http.Request) {
// 	// Set CORS headers if necessary
// 	if config.SetAccessControlHeaders(w, r) {
// 		return // If it's a preflight request, return early.
// 	}
// 	config.SetEnv()

// 	// Initialize Gorilla Mux router
// 	router := mux.NewRouter()

// 	// Define all routes with the corresponding handler functions
// 	router.HandleFunc("/", controller.GetHome).Methods("GET")

// 	// GIS
// 	router.HandleFunc("/data/gis/lokasi", controller.GetRegion).Methods("POST")

// 	// Chat bot inbox
// 	router.HandleFunc("/webhook/nomor/{nomorwa}", controller.PostInboxNomor).Methods("POST")

// 	// Masking list nomor official
// 	router.HandleFunc("/data/phone/all", controller.GetBotList).Methods("GET")

// 	// Helpdesk user data
// 	router.HandleFunc("/data/user/helpdesk/all", controller.GetHelpdeskAll).Methods("GET")
// 	router.HandleFunc("/data/user/helpdesk/masuk", controller.GetLatestHelpdeskMasuk).Methods("GET")
// 	router.HandleFunc("/data/user/helpdesk/selesai", controller.GetLatestHelpdeskSelesai).Methods("GET")

// 	// Pamong desa data
// 	router.HandleFunc("/data/lms/user", controller.GetDataUserFromApi).Methods("GET")

// 	// Testimoni
// 	router.HandleFunc("/data/lms/testi", controller.PostTestimoni).Methods("POST")
// 	router.HandleFunc("/data/lms/random/testi", controller.GetRandomTesti4).Methods("GET")

// 	// Sent items
// 	router.HandleFunc("/data/peserta/sent/{id}", controller.GetSentItem).Methods("GET")

// 	// Unsubscribe feedback
// 	router.HandleFunc("/data/peserta/unsubscribe", controller.PostUnsubscribe).Methods("POST")

// 	// Token linked device
// 	router.HandleFunc("/data/user", controller.PutTokenDataUser).Methods("PUT")

// 	// Add nomor sender for broadcast
// 	router.HandleFunc("/data/sender", controller.PutNomorBlast).Methods("PUT")
// 	router.HandleFunc("/data/sender", controller.GetDataSenders).Methods("GET")

// 	// Blocked senders
// 	router.HandleFunc("/data/blokir", controller.GetDataSendersTerblokir).Methods("GET")

// 	// WA blast recap
// 	router.HandleFunc("/data/rekap", controller.GetRekapBlast).Methods("GET")

// 	// FAQ
// 	router.HandleFunc("/data/faq/{id}", controller.GetFAQ).Methods("GET")

// 	// Legacy
// 	router.HandleFunc("/data/user/task/doing", controller.PutTaskUser).Methods("PUT")
// 	router.HandleFunc("/data/user/task/done", controller.GetTaskDone).Methods("GET")
// 	router.HandleFunc("/data/user/task/done", controller.PostTaskUser).Methods("POST")
// 	router.HandleFunc("/data/pushrepo/kemarin", controller.GetYesterdayDistincWAGroup).Methods("GET")

// 	// Helpdesk
// 	router.HandleFunc("/data/tiket/closed/{id}", controller.GetClosedTicket).Methods("GET")
// 	router.HandleFunc("/data/tiket/rate", controller.PostMasukanTiket).Methods("POST")

// 	// Orders
// 	router.HandleFunc("/data/order/{namalapak}", controller.HandleOrder).Methods("POST")

// 	// User data
// 	router.HandleFunc("/data/user", controller.GetDataUser).Methods("GET")
// 	router.HandleFunc("/auth/register/users", controller.RegisterGmailAuth).Methods("POST")
// 	router.HandleFunc("/data/user", controller.PostDataUser).Methods("POST")
// 	router.HandleFunc("/upload/profpic", controller.UploadProfilePictureHandler).Methods("POST")
// 	router.HandleFunc("/data/user/bio", controller.PostDataBioUser).Methods("POST")

// 	// Projects data
// 	router.HandleFunc("/data/proyek", controller.GetDataProject).Methods("GET")
// 	router.HandleFunc("/data/proyek/approved", controller.GetEditorApprovedProject).Methods("GET")
// 	router.HandleFunc("/data/proyek", controller.PostDataProject).Methods("POST")
// 	router.HandleFunc("/data/metadatabuku", controller.PutMetaDataProject).Methods("PUT")
// 	router.HandleFunc("/data/proyek/publishbuku", controller.PutPublishProject).Methods("PUT")
// 	router.HandleFunc("/data/proyek", controller.PutDataProject).Methods("PUT")
// 	router.HandleFunc("/data/proyek", controller.DeleteDataProject).Methods("DELETE")
// 	router.HandleFunc("/data/proyek/anggota", controller.GetDataMemberProject).Methods("GET")
// 	router.HandleFunc("/data/proyek/editor", controller.GetDataEditorProject).Methods("GET")
// 	router.HandleFunc("/data/proyek/anggota", controller.DeleteDataMemberProject).Methods("DELETE")
// 	router.HandleFunc("/data/proyek/anggota", controller.PostDataMemberProject).Methods("POST")
// 	router.HandleFunc("/data/proyek/editor", controller.PostDataEditorProject).Methods("POST")
// 	router.HandleFunc("/data/proyek/editor", controller.PUtApprovedEditorProject).Methods("PUT")

// 	// File uploads
// 	router.HandleFunc("/upload/coverbuku/{projectid}", controller.UploadCoverBukuWithParamFileHandler).Methods("POST")
// 	router.HandleFunc("/upload/draftbuku/{projectid}", controller.UploadDraftBukuWithParamFileHandler).Methods("POST")
// 	router.HandleFunc("/upload/draftpdfbuku/{projectid}", controller.UploadDraftBukuPDFWithParamFileHandler).Methods("POST")
// 	router.HandleFunc("/upload/sampulpdfbuku/{projectid}", controller.UploadSampulBukuPDFWithParamFileHandler).Methods("POST")
// 	router.HandleFunc("/upload/spk/{projectid}", controller.UploadSPKPDFWithParamFileHandler).Methods("POST")
// 	router.HandleFunc("/upload/spi/{projectid}", controller.UploadSPIPDFWithParamFileHandler).Methods("POST")
// 	router.HandleFunc("/download/draft/{path}", controller.AksesFileRepoDraft).Methods("GET")
// 	router.HandleFunc("/data/proyek/katalog", controller.PostKatalogBuku).Methods("POST")

// 	// More Routes
// 	router.HandleFunc("/notif/ux/postlaporan", controller.PostLaporan).Methods("POST")
// 	router.HandleFunc("/notif/ux/postfeedback", controller.PostFeedback).Methods("POST")
// 	router.HandleFunc("/notif/ux/postmeeting", controller.PostMeeting).Methods("POST")
// 	router.HandleFunc("/notif/ux/postpresensi/{id}", controller.PostPresensi).Methods("POST")
// 	router.HandleFunc("/notif/ux/posttasklists/{id}", controller.PostTaskList).Methods("POST")
// 	router.HandleFunc("/lms/refresh/cookie", controller.RefreshLMSCookie).Methods("GET")
// 	router.HandleFunc("/lms/count/user", controller.GetCountDocUser).Methods("GET")

// 	// Google Auth
// 	router.HandleFunc("/auth/users", controller.Auth).Methods("POST")
// 	router.HandleFunc("/auth/login", controller.GeneratePasswordHandler).Methods("POST")
// 	router.HandleFunc("/auth/verify", controller.VerifyPasswordHandler).Methods("POST")
// 	router.HandleFunc("/auth/resend", controller.ResendPasswordHandler).Methods("POST")

// 	// Products
// 	router.HandleFunc("/products", controller.CreateProduct).Methods("POST")
// 	router.HandleFunc("/products", controller.GetProducts).Methods("GET")
// 	router.HandleFunc("/products/{id}", controller.GetProductByID).Methods("GET")
// 	router.HandleFunc("/products/{id}", controller.UpdateProduct).Methods("PUT")
// 	router.HandleFunc("/products/{id}", controller.DeleteProduct).Methods("DELETE")
// 	router.HandleFunc("/products-export-csv", controller.ExportProductsToCSV).Methods("GET")

// 	// Expense
// 	router.HandleFunc("/expense", controller.CreateExpenseTransaction).Methods("POST")
// 	router.HandleFunc("/expense", controller.GetExpenses).Methods("GET")
// 	router.HandleFunc("/expense/{id}", controller.GetExpenseByID).Methods("GET")
// 	router.HandleFunc("/expense/{id}", controller.UpdateExpense).Methods("PUT")
// 	router.HandleFunc("/expense/{id}", controller.DeleteExpense).Methods("DELETE")
// 	router.HandleFunc("/expense-export-csv", controller.ExportProductsToCSV).Methods("GET")

// 	// Sales
// 	router.HandleFunc("/sales", controller.CreateSalesTransaction).Methods("POST")
// 	router.HandleFunc("/sales", controller.GetSalesTransactions).Methods("GET")
// 	router.HandleFunc("/sales/{id}", controller.GetSalesTransactionByID).Methods("GET")
// 	router.HandleFunc("/sales/{id}", controller.UpdateSalesTransaction).Methods("PUT")
// 	router.HandleFunc("/sales/{id}", controller.DeleteSalesTransaction).Methods("DELETE")
// 	router.HandleFunc("/sales-export-csv", controller.ExportProductsToCSV).Methods("GET")

	
// 	// Customers
// 	router.HandleFunc("/pelanggan", controller.CreateCustomer).Methods("POST")
// 	router.HandleFunc("/pelanggan", controller.GetCustomers).Methods("GET")
// 	router.HandleFunc("/pelanggan/{id}", controller.GetCustomerByID).Methods("GET")
// 	router.HandleFunc("/pelanggan/{id}", controller.UpdateCustomer).Methods("PUT")
// 	router.HandleFunc("/pelanggan/{id}", controller.DeleteCustomer).Methods("DELETE")

// 	// Financial Reports
// 	router.HandleFunc("/laporan", controller.CreateFinancialReport).Methods("POST")
// 	router.HandleFunc("/laporan", controller.GetFinancialReports).Methods("GET")
// 	router.HandleFunc("/laporan/{id}", controller.GetFinancialReportByID).Methods("GET")
// 	router.HandleFunc("/laporan/{id}", controller.DeleteFinancialReport).Methods("DELETE")

// 	// Adding routes to handle dynamic parameters
// 	router.HandleFunc("/data/sale/{id}", controller.GetSalesTransactions).Methods("GET")

// 	// Attach router to default server
// 	http.Handle("/", router)

// 	// Serve HTTP
// 	http.ListenAndServe(":8080", nil)
// }
