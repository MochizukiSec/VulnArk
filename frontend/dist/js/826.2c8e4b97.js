"use strict";(self["webpackChunkvuln_management_frontend"]=self["webpackChunkvuln_management_frontend"]||[]).push([[826],{9868:(e,a,t)=>{var l=t(46518),r=t(79504),n=t(91291),o=t(31240),u=t(72333),i=t(79039),s=RangeError,c=String,d=Math.floor,v=r(u),f=r("".slice),p=r(1..toFixed),b=function(e,a,t){return 0===a?t:a%2===1?b(e,a-1,t*e):b(e*e,a/2,t)},h=function(e){var a=0,t=e;while(t>=4096)a+=12,t/=4096;while(t>=2)a+=1,t/=2;return a},g=function(e,a,t){var l=-1,r=t;while(++l<6)r+=a*e[l],e[l]=r%1e7,r=d(r/1e7)},m=function(e,a){var t=6,l=0;while(--t>=0)l+=e[t],e[t]=d(l/a),l=l%a*1e7},k=function(e){var a=6,t="";while(--a>=0)if(""!==t||0===a||0!==e[a]){var l=c(e[a]);t=""===t?l:t+v("0",7-l.length)+l}return t},y=i((function(){return"0.000"!==p(8e-5,3)||"1"!==p(.9,0)||"1.25"!==p(1.255,2)||"1000000000000000128"!==p(0xde0b6b3a7640080,0)}))||!i((function(){p({})}));l({target:"Number",proto:!0,forced:y},{toFixed:function(e){var a,t,l,r,u=o(this),i=n(e),d=[0,0,0,0,0,0],p="",y="0";if(i<0||i>20)throw new s("Incorrect fraction digits");if(u!==u)return"NaN";if(u<=-1e21||u>=1e21)return c(u);if(u<0&&(p="-",u=-u),u>1e-21)if(a=h(u*b(2,69,1))-69,t=a<0?u*b(2,-a,1):u/b(2,a,1),t*=4503599627370496,a=52-a,a>0){g(d,0,t),l=i;while(l>=7)g(d,1e7,0),l-=7;g(d,b(10,l,1),0),l=a-1;while(l>=23)m(d,1<<23),l-=23;m(d,1<<l),g(d,1,1),m(d,2),y=k(d)}else g(d,0,t),g(d,1<<-a,0),y=k(d)+v("0",i);return i>0?(r=y.length,y=p+(r<=i?"0."+v("0",i-r)+y:f(y,0,r-i)+"."+f(y,r-i))):y=p+y,y}})},20826:(e,a,t)=>{t.r(a),t.d(a,{default:()=>S});t(76918),t(23288),t(9868),t(26099),t(38781);var l=t(20641),r=t(53751),n=t(90033),o={class:"vulndb-container"},u={class:"page-header"},i={class:"header-actions"},s={class:"search-filters"},c={class:"search-box"},d={class:"filter-options"},v={class:"table-header"},f={class:"table-actions"},p={class:"vuln-title"},b={class:"title-text"},h={class:"cvss-score"},g={class:"cvss-value"},m={class:"pagination-container"};function k(e,a,t,k,y,w){var C=(0,l.g2)("el-button"),F=(0,l.g2)("el-input"),_=(0,l.g2)("el-option"),V=(0,l.g2)("el-select"),x=(0,l.g2)("el-card"),D=(0,l.g2)("el-link"),R=(0,l.g2)("el-table-column"),T=(0,l.g2)("el-tag"),S=(0,l.g2)("el-progress"),I=(0,l.g2)("el-table"),L=(0,l.g2)("el-pagination"),W=(0,l.gN)("loading");return(0,l.uX)(),(0,l.CE)("div",o,[(0,l.Lk)("div",u,[a[9]||(a[9]=(0,l.Lk)("div",{class:"header-content"},[(0,l.Lk)("h1",{class:"page-title"},"漏洞库"),(0,l.Lk)("p",{class:"page-subtitle"},"搜索和浏览常见漏洞数据库，支持CVE查询和多种过滤条件")],-1)),(0,l.Lk)("div",i,[(0,l.bF)(C,{type:"success",icon:"el-icon-plus",onClick:k.goToCreate},{default:(0,l.k6)((function(){return a[7]||(a[7]=[(0,l.eW)("新增漏洞")])})),_:1},8,["onClick"]),(0,l.bF)(C,{type:"primary",icon:"el-icon-refresh",onClick:k.fetchVulnerabilities},{default:(0,l.k6)((function(){return a[8]||(a[8]=[(0,l.eW)("刷新")])})),_:1},8,["onClick"])])]),(0,l.bF)(x,{class:"filter-card",shadow:"hover"},{default:(0,l.k6)((function(){return[(0,l.Lk)("div",s,[(0,l.Lk)("div",c,[(0,l.bF)(F,{modelValue:k.searchParams.searchTerm,"onUpdate:modelValue":a[0]||(a[0]=function(e){return k.searchParams.searchTerm=e}),placeholder:"搜索CVE ID、标题或关键词",clearable:"","prefix-icon":"el-icon-search",onKeyup:(0,r.jR)(k.handleSearch,["enter"])},{append:(0,l.k6)((function(){return[(0,l.bF)(C,{onClick:k.handleSearch},{default:(0,l.k6)((function(){return a[10]||(a[10]=[(0,l.eW)("搜索")])})),_:1},8,["onClick"])]})),_:1},8,["modelValue","onKeyup"])]),(0,l.Lk)("div",d,[(0,l.bF)(V,{modelValue:k.searchParams.year,"onUpdate:modelValue":a[1]||(a[1]=function(e){return k.searchParams.year=e}),placeholder:"年份",clearable:"",style:{width:"120px"}},{default:(0,l.k6)((function(){return[((0,l.uX)(!0),(0,l.CE)(l.FK,null,(0,l.pI)(k.yearOptions,(function(e){return(0,l.uX)(),(0,l.Wv)(_,{key:e,label:e?e.toString():"全部",value:e},null,8,["label","value"])})),128))]})),_:1},8,["modelValue"]),(0,l.bF)(V,{modelValue:k.searchParams.severity,"onUpdate:modelValue":a[2]||(a[2]=function(e){return k.searchParams.severity=e}),placeholder:"严重程度",clearable:"",style:{width:"120px"}},{default:(0,l.k6)((function(){return[((0,l.uX)(!0),(0,l.CE)(l.FK,null,(0,l.pI)(k.severityOptions,(function(e){return(0,l.uX)(),(0,l.Wv)(_,{key:e.value,label:e.label,value:e.value},null,8,["label","value"])})),128))]})),_:1},8,["modelValue"]),(0,l.bF)(V,{modelValue:k.searchParams.cvssRange,"onUpdate:modelValue":a[3]||(a[3]=function(e){return k.searchParams.cvssRange=e}),placeholder:"CVSS评分",clearable:"",style:{width:"140px"}},{default:(0,l.k6)((function(){return[((0,l.uX)(!0),(0,l.CE)(l.FK,null,(0,l.pI)(k.cvssRangeOptions,(function(e){return(0,l.uX)(),(0,l.Wv)(_,{key:e.value,label:e.label,value:e.value},null,8,["label","value"])})),128))]})),_:1},8,["modelValue"]),(0,l.bF)(V,{modelValue:k.searchParams.sortBy,"onUpdate:modelValue":a[4]||(a[4]=function(e){return k.searchParams.sortBy=e}),placeholder:"排序字段",style:{width:"140px"}},{default:(0,l.k6)((function(){return[((0,l.uX)(!0),(0,l.CE)(l.FK,null,(0,l.pI)(k.sortOptions,(function(e){return(0,l.uX)(),(0,l.Wv)(_,{key:e.value,label:e.label,value:e.value},null,8,["label","value"])})),128))]})),_:1},8,["modelValue"]),(0,l.bF)(V,{modelValue:k.searchParams.sortOrder,"onUpdate:modelValue":a[5]||(a[5]=function(e){return k.searchParams.sortOrder=e}),placeholder:"排序方式",style:{width:"100px"}},{default:(0,l.k6)((function(){return[(0,l.bF)(_,{key:"asc",label:"升序",value:"asc"}),(0,l.bF)(_,{key:"desc",label:"降序",value:"desc"})]})),_:1},8,["modelValue"]),(0,l.bF)(C,{type:"primary",onClick:k.handleSearch},{default:(0,l.k6)((function(){return a[11]||(a[11]=[(0,l.eW)("应用筛选")])})),_:1},8,["onClick"]),(0,l.bF)(C,{onClick:k.resetFilters},{default:(0,l.k6)((function(){return a[12]||(a[12]=[(0,l.eW)("重置")])})),_:1},8,["onClick"])])])]})),_:1}),(0,l.bF)(x,{shadow:"hover",class:"vulndb-table-card"},{default:(0,l.k6)((function(){return[(0,l.Lk)("div",v,[a[13]||(a[13]=(0,l.Lk)("span",{class:"table-title"},"漏洞列表",-1)),(0,l.Lk)("div",f,[(0,l.bF)(V,{modelValue:k.pageSize,"onUpdate:modelValue":a[6]||(a[6]=function(e){return k.pageSize=e}),placeholder:"每页显示",style:{width:"120px"}},{default:(0,l.k6)((function(){return[(0,l.bF)(_,{value:10,label:"10条/页"}),(0,l.bF)(_,{value:20,label:"20条/页"}),(0,l.bF)(_,{value:50,label:"50条/页"}),(0,l.bF)(_,{value:100,label:"100条/页"})]})),_:1},8,["modelValue"])])]),(0,l.bo)(((0,l.uX)(),(0,l.Wv)(I,{data:k.vulnerabilities,stripe:"",border:"",style:{width:"100%"},class:"vulndb-table",onRowClick:k.handleRowClick},{default:(0,l.k6)((function(){return[(0,l.bF)(R,{prop:"cveId",label:"CVE ID",width:"150",sortable:""},{default:(0,l.k6)((function(e){return[(0,l.bF)(D,{type:"primary",onClick:(0,r.D$)((function(a){return k.goToDetail(e.row.cveId)}),["stop"])},{default:(0,l.k6)((function(){return[(0,l.eW)((0,n.v_)(e.row.cveId),1)]})),_:2},1032,["onClick"])]})),_:1}),(0,l.bF)(R,{prop:"title",label:"漏洞标题","min-width":"200"},{default:(0,l.k6)((function(e){return[(0,l.Lk)("div",p,[(0,l.Lk)("span",b,(0,n.v_)(e.row.title),1),e.row.tags&&e.row.tags.length?((0,l.uX)(),(0,l.Wv)(T,{key:0,size:"small",effect:"plain",class:"tag-item"},{default:(0,l.k6)((function(){return[(0,l.eW)((0,n.v_)(e.row.tags[0]),1)]})),_:2},1024)):(0,l.Q3)("",!0)])]})),_:1}),(0,l.bF)(R,{prop:"cvss",label:"CVSS",width:"100",sortable:""},{default:(0,l.k6)((function(e){return[(0,l.Lk)("div",h,[(0,l.bF)(S,{percentage:10*e.row.cvss,color:k.getCvssColor(e.row.cvss),"show-text":!1,"stroke-width":4,class:"cvss-progress"},null,8,["percentage","color"]),(0,l.Lk)("span",g,(0,n.v_)(e.row.cvss.toFixed(1)),1)])]})),_:1}),(0,l.bF)(R,{prop:"severity",label:"严重程度",width:"120",sortable:""},{default:(0,l.k6)((function(e){return[(0,l.bF)(T,{type:k.getSeverityType(e.row.severity),effect:"dark",size:"small"},{default:(0,l.k6)((function(){return[(0,l.eW)((0,n.v_)(k.getSeverityText(e.row.severity)),1)]})),_:2},1032,["type"])]})),_:1}),(0,l.bF)(R,{prop:"publishedDate",label:"发布日期",width:"120",sortable:""},{default:(0,l.k6)((function(e){return[(0,l.eW)((0,n.v_)(k.formatDate(e.row.publishedDate)),1)]})),_:1}),(0,l.bF)(R,{prop:"lastModifiedDate",label:"更新日期",width:"120",sortable:""},{default:(0,l.k6)((function(e){return[(0,l.eW)((0,n.v_)(k.formatDate(e.row.lastModifiedDate)),1)]})),_:1}),(0,l.bF)(R,{label:"操作",width:"150",fixed:"right"},{default:(0,l.k6)((function(e){return[(0,l.bF)(C,{type:"primary",size:"small",plain:"",icon:"el-icon-view",onClick:(0,r.D$)((function(a){return k.goToDetail(e.row.cveId)}),["stop"])},{default:(0,l.k6)((function(){return a[14]||(a[14]=[(0,l.eW)(" 详情 ")])})),_:2},1032,["onClick"]),(0,l.bF)(C,{type:"success",size:"small",plain:"",icon:"el-icon-plus",onClick:(0,r.D$)((function(a){return k.importToVulnerability(e.row)}),["stop"]),title:"导入到漏洞管理"},{default:(0,l.k6)((function(){return a[15]||(a[15]=[(0,l.eW)(" 导入 ")])})),_:2},1032,["onClick"])]})),_:1})]})),_:1},8,["data","onRowClick"])),[[W,k.loading]]),(0,l.Lk)("div",m,[(0,l.bF)(L,{background:"",layout:"total, prev, pager, next, jumper",total:k.total,"page-size":k.pageSize,"current-page":k.currentPage,onCurrentChange:k.handleCurrentChange},null,8,["total","page-size","current-page","onCurrentChange"])])]})),_:1})])}var y=t(14048),w=t(30388),C=(t(52675),t(89463),t(44114),t(50953)),F=t(75220),_=t(20163),V=t(77918),x=t(94335);const D={name:"VulnDatabaseList",setup:function(){var e=(0,F.rd)(),a=(0,C.KR)(!1),t=(0,C.KR)([]),r=(0,C.KR)(0),n=(0,C.KR)(1),o=(0,C.KR)(20),u=(0,C.Kh)({searchTerm:"",year:null,severity:"",cvssRange:"",sortBy:"publishedDate",sortOrder:"desc"}),i=(0,l.EW)((function(){for(var e=(new Date).getFullYear(),a=[],t=e;t>=2e3;t--)a.push(t);return a})),s=[{value:"critical",label:"严重"},{value:"high",label:"高危"},{value:"medium",label:"中危"},{value:"low",label:"低危"},{value:"info",label:"信息"}],c=[{value:"9-10",label:"9.0 - 10.0"},{value:"7-8.9",label:"7.0 - 8.9"},{value:"4-6.9",label:"4.0 - 6.9"},{value:"0-3.9",label:"0.0 - 3.9"}],d=[{value:"publishedDate",label:"发布日期"},{value:"lastModifiedDate",label:"更新日期"},{value:"cveId",label:"CVE ID"},{value:"cvss",label:"CVSS评分"}],v=function(){var e=(0,w.A)((0,y.A)().mark((function e(){var l,i;return(0,y.A)().wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.prev=0,a.value=!0,l={page:n.value,perPage:o.value,searchTerm:u.searchTerm,severity:u.severity,year:u.year,cvssRange:u.cvssRange,sortBy:u.sortBy,sortOrder:u.sortOrder},console.log("请求漏洞库数据:",{baseURL:x.A.defaults.baseURL,endpoint:"/api/vulndatabase",params:l}),e.next=6,x.A.get("/api/vulndatabase",{params:l});case 6:i=e.sent,console.log("成功获取漏洞库数据:",i.data),t.value=i.data.items||[],r.value=i.data.total||0,0===t.value.length&&r.value>0&&n.value>1&&(n.value=1,v()),e.next=19;break;case 13:e.prev=13,e.t0=e["catch"](0),console.error("获取漏洞库数据失败:",e.t0),_.nk.error("获取漏洞库数据失败，请稍后重试"),t.value=[],r.value=0;case 19:return e.prev=19,a.value=!1,e.finish(19);case 22:case"end":return e.stop()}}),e,null,[[0,13,19,22]])})));return function(){return e.apply(this,arguments)}}(),f=function(){console.log("执行搜索，搜索参数:",u),n.value=1,v()},p=function(){u.searchTerm="",u.year=null,u.severity="",u.cvssRange="",u.sortBy="publishedDate",u.sortOrder="desc",n.value=1,v()},b=function(e){n.value=e,v()},h=function(e){m(e.cveId)},g=function(){e.push("/vulndatabase/create")},m=function(a){e.push("/vulndatabase/".concat(a))},k=function(e){V.s.confirm("确定要将 ".concat(e.cveId," 导入到漏洞管理系统吗？"),"导入确认",{confirmButtonText:"确定",cancelButtonText:"取消",type:"info"}).then((0,w.A)((0,y.A)().mark((function a(){var t;return(0,y.A)().wrap((function(a){while(1)switch(a.prev=a.next){case 0:return a.prev=0,t={title:e.title,description:e.description,cve:e.cveId,cvss:e.cvss,severity:e.severity,remediation:e.remediation||"",references:e.references||[],tags:e.tags||[]},a.next=4,x.A.post("/vulnerabilities/import-from-vulndb",{vulnerability:t});case 4:_.nk.success("成功导入 ".concat(e.cveId," 到漏洞管理系统")),a.next=11;break;case 7:a.prev=7,a.t0=a["catch"](0),console.error("导入漏洞失败:",a.t0),_.nk.error("导入漏洞失败，请稍后重试");case 11:case"end":return a.stop()}}),a,null,[[0,7]])}))))["catch"]((function(){}))},D=function(e){if(!e)return"未知";try{var a=new Date(e);return isNaN(a.getTime())?"日期格式错误":new Intl.DateTimeFormat("zh-CN",{year:"numeric",month:"2-digit",day:"2-digit"}).format(a)}catch(t){return console.error("日期格式化错误:",t,e),"日期格式错误"}},R=function(e){var a={critical:"danger",high:"warning",medium:"warning",low:"success",info:"info"};return a[e]||"info"},T=function(e){var a={critical:"严重",high:"高危",medium:"中危",low:"低危",info:"信息"};return a[e]||"未知"},S=function(e){return e>=9?"#F56C6C":e>=7?"#E6A23C":e>=4?"#F0C050":"#67C23A"};return(0,l.sV)((function(){v()})),{vulnerabilities:t,loading:a,total:r,currentPage:n,pageSize:o,searchParams:u,yearOptions:i,severityOptions:s,cvssRangeOptions:c,sortOptions:d,handleSearch:f,resetFilters:p,fetchVulnerabilities:v,handleCurrentChange:b,handleRowClick:h,goToCreate:g,goToDetail:m,importToVulnerability:k,formatDate:D,getSeverityType:R,getSeverityText:T,getCvssColor:S}}};var R=t(66262);const T=(0,R.A)(D,[["render",k],["__scopeId","data-v-40ea534f"]]),S=T}}]);
//# sourceMappingURL=826.2c8e4b97.js.map