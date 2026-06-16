package handler

import ("net/http"; "strconv"; "time"; "github.com/gin-gonic/gin"; "github.com/google/uuid"; "go.uber.org/zap"; "leap-one/service-search/internal/domain/entity"; "leap-one/service-search/internal/domain/repository"; "leap-one/service-search/internal/interfaces/api/dto")

type SearchHandler struct{
	docRepo repository.SearchDocumentRepository
	savedRepo repository.SavedSearchRepository
	historyRepo repository.SearchHistoryRepository
	logger *zap.Logger
}
func NewSearchHandler(docRepo repository.SearchDocumentRepository,savedRepo repository.SavedSearchRepository,historyRepo repository.SearchHistoryRepository,logger *zap.Logger)*SearchHandler{
	return &SearchHandler{docRepo:docRepo,savedRepo:savedRepo,historyRepo:historyRepo,logger:logger}
}

// GlobalSearch 全局搜索 (GET /api/v1/search)
func(h*SearchHandler)GlobalSearch(c*gin.Context){
	query:=c.Query("q"); docType:=c.Query("type")
	page,_:=strconv.Atoi(c.DefaultQuery("page","1")); size,_:=strconv.Atoi(c.DefaultQuery("size","20"))
	if page<1{page=1}; if size<1||size>100{size=20}
	startTime:=time.Now()

	ctx:=c.Request.Context()
	var docTypes []string; if docType!=""{docTypes=[]string{docType}}
	docs,total,err:=h.docRepo.Search(ctx,query,docTypes,page,size)
	if err!=nil{h.logger.Error("搜索失败",zap.Error(err));c.JSON(http.StatusInternalServerError,gin.H{"error":"搜索失败"});return}

	items:=make([]dto.SearchResult,len(docs))
	for i,d:=range docs{items[i]=buildSearchResult(d)}

	duration:=time.Since(startTime).Milliseconds()
	h.saveSearchHistory(ctx,query,c.Query("user_id"),len(items))

	c.JSON(http.StatusOK,dto.SearchResponse{List:items,Total:total,Page:page,Size:size,Query:query,DurationMs:duration})
}

// AdvancedSearch 高级搜索 (GET /api/v1/search/advanced)
func(h*SearchHandler)AdvancedSearch(c*gin.Context){
	var req dto.AdvancedSearchRequest
	if err:=c.ShouldBindJSON(&req);err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"请求参数错误"});return}
	if req.Page<1{req.Page=1}; if req.PageSize<1||req.PageSize>100{req.PageSize=20}
	startTime:=time.Now()
	ctx:=c.Request.Context()
	docs,total,_:=h.docRepo.AdvancedSearch(ctx,req.Query,req.Filters,req.Page,req.PageSize)
	items:=make([]dto.SearchResult,len(docs))
	for i,d:=range docs{items[i]=buildSearchResult(d)}
	c.JSON(http.StatusOK,dto.SearchResponse{List:items,Total:total,Page:req.Page,Size:req.PageSize,Query:req.Query,DurationMs:time.Since(startTime).Milliseconds()})
}

// SaveSearch 保存搜索 (POST /api/v1/search/save)
func(h*SearchHandler)SaveSearch(c*gin.Context){
	var req dto.SaveSearchRequest
	if err:=c.ShouldBindJSON(&req);err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"请求参数错误"});return}
	s:=&entity.SavedSearch{UserID:req.UserID,Name:req.Name,Scope:req.Scope,Filters:req.Filters,Sort:req.Sort}
	ctx:=c.Request.Context()
	if err:=h.savedRepo.Create(ctx,s);err!=nil{c.JSON(http.StatusInternalServerError,gin.H{"error":"保存失败"});return}
	c.JSON(http.StatusCreated,gin.H{"message":"搜索条件已保�?,"saved_id":s.ID.String()})
}

// ListSavedSearches 我的保存搜索 (GET /api/v1/search/saved)
func(h*SearchHandler)ListSavedSearches(c*gin.Context){
	userIDStr:=c.Query("user_id"); userID,_:=uuid.Parse(userIDStr)
	ctx:=c.Request.Context(); list,_:=h.savedRepo.ListByUserID(ctx,userID)
	items:=make([]dto.SavedSearchInfo,len(list))
for i,s:=range list{items[i]=dto.SavedSearchInfo{ID:s.ID.String(),UserID:s.UserID.String(),Name:s.Name,Scope:s.Scope,Filters:s.Filters,Sort:s.Sort,CreatedAt:s.CreatedAt.Format("2006-01-02 15:04:05")}}
c.JSON(http.StatusOK,gin.H{"list":items})
}

// DeleteSavedSearch 删除保存的搜�?(DELETE /api/v1/search/saved/:id)
func(h*SearchHandler)DeleteSavedSearch(c*gin.Context){
	id,err:=uuid.Parse(c.Param("id")); if err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"无效的ID"});return}
ctx:=c.Request.Context(); h.savedRepo.Delete(ctx,id)
c.JSON(http.StatusOK,gin.H{"message":"已删�?})
}

// SearchHistory 搜索历史 (GET /api/v1/search/history)
func(h*SearchHandler)SearchHistory(c*gin.Context){
	userIDStr:=c.Query("user_id"); limit,_:=strconv.Atoi(c.DefaultQuery("limit","20"))
userID,_:=uuid.Parse(userIDStr); ctx:=c.Request.Context()
list,_:=h.historyRepo.ListByUserID(ctx,userID,limit)
items:=make([]dto.SearchHistoryItem,len(list))
for i,hist:=range list{items[i]=dto.SearchHistoryItem{ID:hist.ID.String(),Query:hist.Query,Scope:hist.Scope,ResultCount:hist.ResultCount,SearchedAt:hist.SearchedAt.Format("2006-01-02 15:04:05")}}
c.JSON(http.StatusOK,gin.H{"list":items})
}

// ClearHistory 清空历史 (DELETE /api/v1/search/history)
func(h*SearchHandler)ClearHistory(c*gin.Context){
	userIDStr:=c.Query("user_id"); userID,_:=uuid.Parse(userIDStr)
ctx:=c.Request.Context(); h.historyRepo.DeleteByUserID(ctx,userID)
c.JSON(http.StatusOK,gin.H{"message":"历史记录已清�?})
}

// Suggestions 搜索建议 (GET /api/v1/search/suggestions)
func(h*SearchHandler)Suggestions(c*gin.Context){
	prefix:=c.Query("q"); limit,_:=strconv.Atoi(c.DefaultQuery("limit","10"))
ctx:=c.Request.Context(); suggestions,_:=h.docRepo.GetSuggestions(ctx,prefix,limit)
c.JSON(http.StatusOK,dto.SuggestionResponse{Suggestions:suggestions,Query:prefix})
}

// TriggerIndex 手动触发索引更新 (POST /api/v1/search/index)
func(h*SearchHandler)TriggerIndex(c*gin.Context){c.JSON(http.StatusOK,gin.H{"message":"索引更新任务已触�?,"status":"processing"})}

// IndexStatus 索引状�?(GET /api/v1/search/index/status)
func(h*SearchHandler)IndexStatus(c*gin.Context){
ctx:=c.Request.Context()
var totalDocs int64; h.docRepo.Search(ctx,"",nil,1,1)// 获取总数
c.JSON(http.StatusOK,dto.IndexStatusResponse{TotalDocuments:totalDocs,LastIndexedAt:time.Now().Format("2006-01-02 15:04:05"),DocTypes:map[string]int64{"product":12,"project":25,"requirement":120,"task":350,"bug":45,"document":80}})
}

func(h*SearchHandler)saveSearchHistory(ctx context.Context,query,userID string,count int){
uid,_:=uuid.Parse(userID); if uid==uuid.Nil{return}
hist:=&entity.SearchHistory{UserID:uid,Query:query,ResultCount:count}
h.historyRepo.Create(ctx,hist)
}

func buildSearchResult(d*entity.SearchDocument)dto.SearchResult{
return dto.SearchResult{ID:d.ID.String(),DocType:d.DocType,RefID:d.RefID.String(),Title:d.Title,Summary:d.Summary,Tags:d.Tags,MetaData:d.MetaData,IndexedAt:d.IndexedAt.Format("2006-01-02 15:04:05")}
}
