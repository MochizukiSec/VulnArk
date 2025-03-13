"use strict";(self["webpackChunkvuln_management_frontend"]=self["webpackChunkvuln_management_frontend"]||[]).push([[467],{9868:(e,t,n)=>{var a=n(46518),r=n(79504),l=n(91291),i=n(31240),s=n(72333),c=n(79039),o=RangeError,u=String,d=Math.floor,f=r(s),p=r("".slice),v=r(1..toFixed),k=function(e,t,n){return 0===t?n:t%2===1?k(e,t-1,n*e):k(e*e,t/2,n)},h=function(e){var t=0,n=e;while(n>=4096)t+=12,n/=4096;while(n>=2)t+=1,n/=2;return t},b=function(e,t,n){var a=-1,r=n;while(++a<6)r+=t*e[a],e[a]=r%1e7,r=d(r/1e7)},g=function(e,t){var n=6,a=0;while(--n>=0)a+=e[n],e[n]=d(a/t),a=a%t*1e7},m=function(e){var t=6,n="";while(--t>=0)if(""!==n||0===t||0!==e[t]){var a=u(e[t]);n=""===n?a:n+f("0",7-a.length)+a}return n},y=c((function(){return"0.000"!==v(8e-5,3)||"1"!==v(.9,0)||"1.25"!==v(1.255,2)||"1000000000000000128"!==v(0xde0b6b3a7640080,0)}))||!c((function(){v({})}));a({target:"Number",proto:!0,forced:y},{toFixed:function(e){var t,n,a,r,s=i(this),c=l(e),d=[0,0,0,0,0,0],v="",y="0";if(c<0||c>20)throw new o("Incorrect fraction digits");if(s!==s)return"NaN";if(s<=-1e21||s>=1e21)return u(s);if(s<0&&(v="-",s=-s),s>1e-21)if(t=h(s*k(2,69,1))-69,n=t<0?s*k(2,-t,1):s/k(2,t,1),n*=4503599627370496,t=52-t,t>0){b(d,0,n),a=c;while(a>=7)b(d,1e7,0),a-=7;b(d,k(10,a,1),0),a=t-1;while(a>=23)g(d,1<<23),a-=23;g(d,1<<a),b(d,1,1),g(d,2),y=m(d)}else b(d,0,n),b(d,1<<-t,0),y=m(d)+f("0",c);return c>0?(r=y.length,y=v+(r<=c?"0."+f("0",c-r)+y:p(y,0,r-c)+"."+p(y,r-c))):y=v+y,y}})},17467:(e,t,n)=>{n.r(t),n.d(t,{default:()=>M});n(9868),n(27495),n(5746);var a=n(20641),r=n(53751),l=n(90033),i={class:"vulnerability-list"},s={class:"page-header"},c={class:"header-actions"},o={class:"filter-card"},u={class:"option-with-tag"},d={class:"option-with-tag"},f={class:"option-with-tag"},p={class:"option-with-tag"},v={class:"option-with-tag"},k={class:"option-with-tag"},h={class:"option-with-tag"},b={class:"option-with-tag"},g={class:"option-with-tag"},m={class:"option-with-tag"},y={class:"vuln-table-card"},w={key:0,class:"search-tags"},_={class:"vuln-title"},C={key:0,class:"el-icon-circle-check"},F={key:1,class:"el-icon-loading"},L={key:2,class:"el-icon-warning"},W={key:3,class:"el-icon-close"},S={key:4,class:"el-icon-question"},T={key:0,class:"cve-id"},P={key:1,class:"cve-id empty"},z={class:"cvss-container"},E={key:1,class:"cvss-score empty"},V={class:"date-text"},x={class:"action-buttons"},A={key:1,class:"no-data"},X={class:"pagination-container"};function U(e,t,n,U,q,B){var D=(0,a.g2)("el-button"),O=(0,a.g2)("el-tag"),Q=(0,a.g2)("el-option"),$=(0,a.g2)("el-select"),j=(0,a.g2)("el-form-item"),I=(0,a.g2)("el-input"),N=(0,a.g2)("el-form"),K=(0,a.g2)("router-link"),R=(0,a.g2)("el-table-column"),M=(0,a.g2)("el-dropdown-item"),G=(0,a.g2)("el-dropdown-menu"),H=(0,a.g2)("el-dropdown"),J=(0,a.g2)("el-table"),Y=(0,a.g2)("el-pagination"),Z=(0,a.gN)("loading");return(0,a.uX)(),(0,a.CE)("div",i,[(0,a.Lk)("div",s,[t[5]||(t[5]=(0,a.Fv)('<div class="header-content" data-v-1dab35a1><h1 class="page-title" data-v-1dab35a1><span class="title-icon" data-v-1dab35a1><i class="el-icon-warning" data-v-1dab35a1></i></span> 漏洞列表 <span class="title-highlight" data-v-1dab35a1>管理</span></h1><p class="page-subtitle" data-v-1dab35a1>查看和管理所有漏洞信息，监控安全风险</p></div>',1)),(0,a.Lk)("div",c,[(0,a.bF)(D,{type:"primary",onClick:U.goToCreate,class:"action-btn"},{default:(0,a.k6)((function(){return t[3]||(t[3]=[(0,a.Lk)("i",{class:"el-icon-plus"},null,-1),(0,a.eW)(" 添加漏洞 ")])})),_:1},8,["onClick"]),(0,a.bF)(D,{onClick:U.goToImport,class:"action-btn"},{default:(0,a.k6)((function(){return t[4]||(t[4]=[(0,a.Lk)("i",{class:"el-icon-upload2"},null,-1),(0,a.eW)(" 导入漏洞 ")])})),_:1},8,["onClick"])])]),(0,a.Lk)("div",o,[t[18]||(t[18]=(0,a.Lk)("div",{class:"filter-header"},[(0,a.Lk)("i",{class:"el-icon-search"}),(0,a.Lk)("span",null,"搜索条件")],-1)),(0,a.bF)(N,{inline:!0,model:U.searchParams,class:"search-form"},{default:(0,a.k6)((function(){return[(0,a.bF)(j,{label:"严重程度"},{default:(0,a.k6)((function(){return[(0,a.bF)($,{modelValue:U.searchParams.severity,"onUpdate:modelValue":t[0]||(t[0]=function(e){return U.searchParams.severity=e}),placeholder:"选择严重程度",clearable:"","popper-class":"wider-dropdown"},{default:(0,a.k6)((function(){return[(0,a.bF)(Q,{label:"严重",value:"critical"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",u,[(0,a.bF)(O,{type:"danger",effect:"dark",size:"small"},{default:(0,a.k6)((function(){return t[6]||(t[6]=[(0,a.eW)("严重")])})),_:1})])]})),_:1}),(0,a.bF)(Q,{label:"高危",value:"high"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",d,[(0,a.bF)(O,{type:"warning",effect:"dark",size:"small"},{default:(0,a.k6)((function(){return t[7]||(t[7]=[(0,a.eW)("高危")])})),_:1})])]})),_:1}),(0,a.bF)(Q,{label:"中危",value:"medium"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",f,[(0,a.bF)(O,{type:"info",effect:"dark",size:"small"},{default:(0,a.k6)((function(){return t[8]||(t[8]=[(0,a.eW)("中危")])})),_:1})])]})),_:1}),(0,a.bF)(Q,{label:"低危",value:"low"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",p,[(0,a.bF)(O,{type:"success",effect:"dark",size:"small"},{default:(0,a.k6)((function(){return t[9]||(t[9]=[(0,a.eW)("低危")])})),_:1})])]})),_:1}),(0,a.bF)(Q,{label:"信息",value:"info"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",v,[(0,a.bF)(O,{type:"info",effect:"plain",size:"small"},{default:(0,a.k6)((function(){return t[10]||(t[10]=[(0,a.eW)("信息")])})),_:1})])]})),_:1})]})),_:1},8,["modelValue"])]})),_:1}),(0,a.bF)(j,{label:"状态"},{default:(0,a.k6)((function(){return[(0,a.bF)($,{modelValue:U.searchParams.status,"onUpdate:modelValue":t[1]||(t[1]=function(e){return U.searchParams.status=e}),placeholder:"选择状态",clearable:"","popper-class":"wider-dropdown"},{default:(0,a.k6)((function(){return[(0,a.bF)(Q,{label:"开放",value:"open"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",k,[(0,a.bF)(O,{type:"danger",effect:"plain",size:"small"},{default:(0,a.k6)((function(){return t[11]||(t[11]=[(0,a.eW)("开放")])})),_:1})])]})),_:1}),(0,a.bF)(Q,{label:"处理中",value:"in_progress"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",h,[(0,a.bF)(O,{type:"warning",effect:"plain",size:"small"},{default:(0,a.k6)((function(){return t[12]||(t[12]=[(0,a.eW)("处理中")])})),_:1})])]})),_:1}),(0,a.bF)(Q,{label:"已解决",value:"resolved"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",b,[(0,a.bF)(O,{type:"success",effect:"plain",size:"small"},{default:(0,a.k6)((function(){return t[13]||(t[13]=[(0,a.eW)("已解决")])})),_:1})])]})),_:1}),(0,a.bF)(Q,{label:"已关闭",value:"closed"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",g,[(0,a.bF)(O,{type:"info",effect:"plain",size:"small"},{default:(0,a.k6)((function(){return t[14]||(t[14]=[(0,a.eW)("已关闭")])})),_:1})])]})),_:1}),(0,a.bF)(Q,{label:"误报",value:"false_positive"},{default:(0,a.k6)((function(){return[(0,a.Lk)("div",m,[(0,a.bF)(O,{effect:"plain",size:"small"},{default:(0,a.k6)((function(){return t[15]||(t[15]=[(0,a.eW)("误报")])})),_:1})])]})),_:1})]})),_:1},8,["modelValue"])]})),_:1}),(0,a.bF)(j,{label:"搜索"},{default:(0,a.k6)((function(){return[(0,a.bF)(I,{modelValue:U.searchParams.searchTerm,"onUpdate:modelValue":t[2]||(t[2]=function(e){return U.searchParams.searchTerm=e}),placeholder:"搜索标题、描述或CVE编号",clearable:"","prefix-icon":"el-icon-search",onKeyup:(0,r.jR)(U.search,["enter"]),class:"search-input"},null,8,["modelValue","onKeyup"])]})),_:1}),(0,a.bF)(j,null,{default:(0,a.k6)((function(){return[(0,a.bF)(D,{type:"primary",onClick:U.search,class:"search-btn"},{default:(0,a.k6)((function(){return t[16]||(t[16]=[(0,a.Lk)("i",{class:"el-icon-search"},null,-1),(0,a.eW)(" 搜索 ")])})),_:1},8,["onClick"]),(0,a.bF)(D,{onClick:U.resetSearch,class:"reset-btn"},{default:(0,a.k6)((function(){return t[17]||(t[17]=[(0,a.Lk)("i",{class:"el-icon-refresh"},null,-1),(0,a.eW)(" 重置 ")])})),_:1},8,["onClick"])]})),_:1})]})),_:1},8,["model"])]),(0,a.Lk)("div",y,[U.hasActiveFilters?((0,a.uX)(),(0,a.CE)("div",w,[t[23]||(t[23]=(0,a.Lk)("span",{class:"search-tags-label"},"当前筛选条件:",-1)),U.searchParams.searchTerm?((0,a.uX)(),(0,a.Wv)(O,{key:0,closable:"",onClose:U.clearSearchTerm,type:"info",effect:"plain",class:"filter-tag"},{default:(0,a.k6)((function(){return[t[19]||(t[19]=(0,a.Lk)("i",{class:"el-icon-search"},null,-1)),(0,a.eW)(" 关键词: "+(0,l.v_)(U.searchParams.searchTerm),1)]})),_:1},8,["onClose"])):(0,a.Q3)("",!0),U.searchParams.severity?((0,a.uX)(),(0,a.Wv)(O,{key:1,closable:"",onClose:U.clearSeverity,type:U.getSeverityType(U.searchParams.severity),effect:"dark",class:"filter-tag"},{default:(0,a.k6)((function(){return[t[20]||(t[20]=(0,a.Lk)("i",{class:"el-icon-warning-outline"},null,-1)),(0,a.eW)(" 严重程度: "+(0,l.v_)(U.getSeverityLabel(U.searchParams.severity)),1)]})),_:1},8,["onClose","type"])):(0,a.Q3)("",!0),U.searchParams.status?((0,a.uX)(),(0,a.Wv)(O,{key:2,closable:"",onClose:U.clearStatus,type:U.getStatusType(U.searchParams.status),effect:"plain",class:"filter-tag"},{default:(0,a.k6)((function(){return[t[21]||(t[21]=(0,a.Lk)("i",{class:"el-icon-info"},null,-1)),(0,a.eW)(" 状态: "+(0,l.v_)(U.getStatusText(U.searchParams.status)),1)]})),_:1},8,["onClose","type"])):(0,a.Q3)("",!0),(0,a.bF)(D,{size:"small",type:"primary",plain:"",onClick:U.resetSearch,class:"clear-all-btn"},{default:(0,a.k6)((function(){return t[22]||(t[22]=[(0,a.Lk)("i",{class:"el-icon-delete"},null,-1),(0,a.eW)(" 清除所有筛选 ")])})),_:1},8,["onClick"])])):(0,a.Q3)("",!0),(0,a.bo)(((0,a.uX)(),(0,a.Wv)(J,{data:U.vulnerabilities,style:{width:"100%"},onSortChange:U.handleSortChange,border:"",stripe:"","highlight-current-row":"",class:"vuln-table"},{default:(0,a.k6)((function(){return[(0,a.bF)(R,{prop:"title",label:"标题","min-width":"200",sortable:"custom"},{default:(0,a.k6)((function(e){return[(0,a.bF)(K,{to:"/vulnerabilities/".concat(e.row.id),class:"vuln-link"},{default:(0,a.k6)((function(){return[(0,a.Lk)("span",_,(0,l.v_)(e.row.title),1)]})),_:2},1032,["to"])]})),_:1}),(0,a.bF)(R,{prop:"severity",label:"严重程度",width:"120",sortable:"custom"},{default:(0,a.k6)((function(n){return[(0,a.bF)(O,{type:e.$filters.severityClass(n.row.severity),effect:"dark",size:"small",class:"severity-tag"},{default:(0,a.k6)((function(){return[t[24]||(t[24]=(0,a.Lk)("i",{class:"el-icon-warning-outline severity-icon"},null,-1)),(0,a.eW)(" "+(0,l.v_)(e.$filters.severityText(n.row.severity)),1)]})),_:2},1032,["type"])]})),_:1}),(0,a.bF)(R,{prop:"status",label:"状态",width:"120",sortable:"custom"},{default:(0,a.k6)((function(t){return[(0,a.bF)(O,{type:e.$filters.statusClass(t.row.status),effect:"plain",size:"small",class:"status-tag"},{default:(0,a.k6)((function(){return["resolved"===t.row.status?((0,a.uX)(),(0,a.CE)("i",C)):"in_progress"===t.row.status?((0,a.uX)(),(0,a.CE)("i",F)):"open"===t.row.status?((0,a.uX)(),(0,a.CE)("i",L)):"closed"===t.row.status?((0,a.uX)(),(0,a.CE)("i",W)):((0,a.uX)(),(0,a.CE)("i",S)),(0,a.eW)(" "+(0,l.v_)(e.$filters.statusText(t.row.status)),1)]})),_:2},1032,["type"])]})),_:1}),(0,a.bF)(R,{prop:"cve",label:"CVE编号",width:"150"},{default:(0,a.k6)((function(e){return[e.row.cve?((0,a.uX)(),(0,a.CE)("span",T,[t[25]||(t[25]=(0,a.Lk)("i",{class:"el-icon-document cve-icon"},null,-1)),(0,a.eW)(" "+(0,l.v_)(e.row.cve),1)])):((0,a.uX)(),(0,a.CE)("span",P,"-"))]})),_:1}),(0,a.bF)(R,{prop:"cvss",label:"CVSS评分",width:"120",sortable:"custom"},{default:(0,a.k6)((function(e){return[(0,a.Lk)("div",z,[e.row.cvss?((0,a.uX)(),(0,a.CE)("div",{key:0,class:(0,l.C4)(["cvss-score",U.getCvssClass(e.row.cvss)])},(0,l.v_)(e.row.cvss.toFixed(1)),3)):((0,a.uX)(),(0,a.CE)("div",E,"-"))])]})),_:1}),(0,a.bF)(R,{prop:"discoveredAt",label:"发现时间",width:"180",sortable:"custom"},{default:(0,a.k6)((function(n){return[(0,a.Lk)("span",V,[t[26]||(t[26]=(0,a.Lk)("i",{class:"el-icon-time date-icon"},null,-1)),(0,a.eW)(" "+(0,l.v_)(e.$filters.formatDate(n.row.discoveredAt)),1)])]})),_:1}),(0,a.bF)(R,{label:"操作",width:"180",fixed:"right"},{default:(0,a.k6)((function(e){return[(0,a.Lk)("div",x,[(0,a.bF)(D,{size:"small",type:"primary",icon:"el-icon-view",circle:"",onClick:function(t){return U.viewDetails(e.row.id)},title:"查看详情"},null,8,["onClick"]),(0,a.bF)(H,{trigger:"click",size:"small",placement:"bottom-end"},{dropdown:(0,a.k6)((function(){return[(0,a.bF)(G,null,{default:(0,a.k6)((function(){return[(0,a.bF)(M,{disabled:"open"===e.row.status,onClick:function(t){return U.quickUpdateStatus(e.row.id,"open")}},{default:(0,a.k6)((function(){return t[27]||(t[27]=[(0,a.Lk)("i",{class:"el-icon-warning status-icon"},null,-1),(0,a.eW)(" 开放 ")])})),_:2},1032,["disabled","onClick"]),(0,a.bF)(M,{disabled:"in_progress"===e.row.status,onClick:function(t){return U.quickUpdateStatus(e.row.id,"in_progress")}},{default:(0,a.k6)((function(){return t[28]||(t[28]=[(0,a.Lk)("i",{class:"el-icon-loading status-icon"},null,-1),(0,a.eW)(" 处理中 ")])})),_:2},1032,["disabled","onClick"]),(0,a.bF)(M,{disabled:"resolved"===e.row.status,onClick:function(t){return U.quickUpdateStatus(e.row.id,"resolved")}},{default:(0,a.k6)((function(){return t[29]||(t[29]=[(0,a.Lk)("i",{class:"el-icon-circle-check status-icon"},null,-1),(0,a.eW)(" 已解决 ")])})),_:2},1032,["disabled","onClick"]),(0,a.bF)(M,{disabled:"closed"===e.row.status,onClick:function(t){return U.quickUpdateStatus(e.row.id,"closed")}},{default:(0,a.k6)((function(){return t[30]||(t[30]=[(0,a.Lk)("i",{class:"el-icon-close status-icon"},null,-1),(0,a.eW)(" 已关闭 ")])})),_:2},1032,["disabled","onClick"]),(0,a.bF)(M,{disabled:"false_positive"===e.row.status,onClick:function(t){return U.quickUpdateStatus(e.row.id,"false_positive")}},{default:(0,a.k6)((function(){return t[31]||(t[31]=[(0,a.Lk)("i",{class:"el-icon-question status-icon"},null,-1),(0,a.eW)(" 误报 ")])})),_:2},1032,["disabled","onClick"])]})),_:2},1024)]})),default:(0,a.k6)((function(){return[(0,a.bF)(D,{size:"small",type:"warning",icon:"el-icon-edit",circle:"",title:"更改状态"})]})),_:2},1024),(0,a.bF)(D,{size:"small",type:"danger",icon:"el-icon-delete",circle:"",onClick:function(t){return U.confirmDelete(e.row)},title:"删除漏洞"},null,8,["onClick"])])]})),_:1})]})),_:1},8,["data","onSortChange"])),[[Z,U.loading]]),!U.loading&&U.vulnerabilities&&0===U.vulnerabilities.length?((0,a.uX)(),(0,a.CE)("div",A,[t[33]||(t[33]=(0,a.Lk)("i",{class:"el-icon-document"},null,-1)),t[34]||(t[34]=(0,a.Lk)("p",null,"暂无漏洞数据",-1)),(0,a.bF)(D,{type:"primary",size:"small",onClick:U.goToCreate},{default:(0,a.k6)((function(){return t[32]||(t[32]=[(0,a.eW)(" 创建漏洞 ")])})),_:1},8,["onClick"])])):(0,a.Q3)("",!0),(0,a.Lk)("div",X,[(0,a.bF)(Y,{background:"","current-page":U.currentPage,"page-sizes":[10,20,50,100],"page-size":U.perPage,layout:"total, sizes, prev, pager, next, jumper",total:U.totalVulnerabilities,onSizeChange:U.handleSizeChange,onCurrentChange:U.handleCurrentChange},null,8,["current-page","page-size","total","onSizeChange","onCurrentChange"])])])])}var q=n(14048),B=n(30388),D=n(29201),O=(n(44114),n(50953)),Q=n(66278),$=n(75220),j=n(20163),I=n(77918);const N={name:"VulnerabilityList",setup:function(){var e=(0,Q.Pj)(),t=(0,$.rd)(),n=(0,a.EW)((function(){return e.getters["vulnerability/isLoading"]})),r=(0,O.Kh)({severity:"",status:"",assignedTo:"",searchTerm:"",sortBy:"createdAt",sortOrder:"desc",page:1,perPage:10}),l=(0,a.EW)((function(){return r.searchTerm||r.severity||r.status})),i=(0,a.EW)((function(){return e.getters["vulnerability/allVulnerabilities"]||[]})),s=(0,a.EW)((function(){return e.getters["vulnerability/pagination"]||{total:0}})),c=(0,a.EW)((function(){return r.page})),o=(0,a.EW)((function(){return r.perPage})),u=(0,a.EW)((function(){return s.value&&s.value.total||0})),d=function(){r.searchTerm="",b()},f=function(){r.severity="",b()},p=function(){r.status="",b()},v=function(e){var t={critical:"严重",high:"高危",medium:"中危",low:"低危",info:"信息"};return t[e]||e},k=function(e){var t={critical:"danger",high:"warning",medium:"",low:"success",info:"info"};return t[e]||""},h=function(e){var t={open:"danger",in_progress:"warning",resolved:"success",closed:"info",false_positive:""};return t[e]||""},b=function(){console.log("执行搜索，搜索参数:",r);var t={searchTerm:r.searchTerm||"无",severity:r.severity||"无",status:r.status||"无"};if(console.log("搜索条件摘要:",t),r.searchTerm&&j.nk.info('正在搜索："'.concat(r.searchTerm,'"...')),r.severity){var n={critical:"严重",high:"高危",medium:"中危",low:"低危",info:"信息"};j.nk.info("筛选严重程度：".concat(n[r.severity]))}if(r.status){var a={open:"开放",in_progress:"处理中",resolved:"已解决",closed:"已关闭",false_positive:"误报"};j.nk.info("筛选状态：".concat(a[r.status]))}e.dispatch("vulnerability/updateSearchParams",(0,D.A)({},r))},g=function(){Object.assign(r,{severity:"",status:"",assignedTo:"",searchTerm:"",sortBy:"createdAt",sortOrder:"desc",page:1}),console.log("重置搜索条件:",r),b()},m=function(t){var n=t.prop,a=t.order;if(n){var r={title:"title",severity:"severity",status:"status",cvss:"cvss",discoveredAt:"discoveredAt",createdAt:"createdAt"},l=r[n]||"createdAt",i="ascending"===a?"asc":"desc";console.log("表格排序变化:",{prop:n,order:a,sortBy:l,sortOrder:i}),e.dispatch("vulnerability/updateSearchParams",{sortBy:l,sortOrder:i})}},y=function(t){console.log("每页显示数量变更为:",t),e.dispatch("vulnerability/updateSearchParams",{perPage:t})},w=function(t){console.log("当前页码变更为:",t),e.dispatch("vulnerability/updateSearchParams",{page:t})},_=function(e){t.push("/vulnerabilities/".concat(e))},C=function(){t.push("/vulnerabilities/create")},F=function(){t.push("/vulnerabilities/import")},L=function(){var t=(0,B.A)((0,q.A)().mark((function t(n,a){return(0,q.A)().wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.prev=0,t.next=3,e.dispatch("vulnerability/updateVulnerability",{id:n,data:{status:a}});case 3:j.nk.success("漏洞状态已更新为".concat(W(a))),e.dispatch("vulnerability/fetchVulnerabilities"),t.next=10;break;case 7:t.prev=7,t.t0=t["catch"](0),j.nk.error("更新状态失败: ".concat(t.t0.message||"未知错误"));case 10:case"end":return t.stop()}}),t,null,[[0,7]])})));return function(e,n){return t.apply(this,arguments)}}(),W=function(e){var t={open:"开放",in_progress:"处理中",resolved:"已解决",closed:"已关闭",false_positive:"误报"};return t[e]||e},S=function(e){I.s.confirm('确定要删除漏洞 "'.concat(e.title,'" 吗？'),"警告",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then((function(){T(e.id)}))["catch"]((function(){(0,j.nk)({type:"info",message:"已取消删除"})}))},T=function(){var t=(0,B.A)((0,q.A)().mark((function t(n){return(0,q.A)().wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.prev=0,t.next=3,e.dispatch("vulnerability/deleteVulnerability",n);case 3:e.dispatch("dashboard/fetchDashboardData"),(0,j.nk)({type:"success",message:"删除成功"}),t.next=10;break;case 7:t.prev=7,t.t0=t["catch"](0),(0,j.nk)({type:"error",message:"删除失败: "+(t.t0.message||"未知错误")});case 10:case"end":return t.stop()}}),t,null,[[0,7]])})));return function(e){return t.apply(this,arguments)}}(),P=function(e){return e?e>=9?"critical":e>=7?"high":e>=4?"medium":e>=.1?"low":"":""};return(0,a.wB)(r,(function(e){console.log("搜索参数已更新:",e)}),{deep:!0}),(0,a.sV)((function(){console.log("VulnerabilityList组件已挂载，准备获取数据"),g(),e.dispatch("vulnerability/fetchVulnerabilities").then((function(e){console.log("初始漏洞数据加载成功:",e)}))["catch"]((function(e){console.error("初始漏洞数据加载失败:",e),j.nk.error("加载漏洞列表失败，请刷新页面重试")}))})),{loading:n,searchParams:r,vulnerabilities:i,pagination:s,currentPage:c,perPage:o,totalVulnerabilities:u,hasActiveFilters:l,search:b,resetSearch:g,clearSearchTerm:d,clearSeverity:f,clearStatus:p,getSeverityLabel:v,getSeverityType:k,getStatusType:h,handleSortChange:m,handleSizeChange:y,handleCurrentChange:w,viewDetails:_,goToCreate:C,goToImport:F,confirmDelete:S,getCvssClass:P,quickUpdateStatus:L,getStatusText:W}}};var K=n(66262);const R=(0,K.A)(N,[["render",U],["__scopeId","data-v-1dab35a1"]]),M=R}}]);
//# sourceMappingURL=467.e4ba1570.js.map