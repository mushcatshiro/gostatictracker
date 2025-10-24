package render

var bkmkSetupElmGroup = []elm{
	{tag: "h2", innerText: "Your Bookmarklet"},
	{tag: "p", innerText: "Drag this link to your bookmarks bar:"},
}

const bkmkStyleString = `padding: 10px 20px;
	background-color: #4CAF50;
	color: white;
	text-decoration: none;
	border-radius: 4px;`

const bkmkScript = `javascript:void((function(){
	function getMetaValue(propName){
		const metas=document.getElementsByTagName('meta');
		for(let i=0;i<metas.length;i++){
			const metaName=metas[i].getAttribute('name')||metas[i].getAttribute('property');
			if(metaName===propName){
				return metas[i].getAttribute('content');
			}
		}
		return'';
	}
	const metaDescription=getMetaValue('og:description')||getMetaValue('description')||'';
	window.open('%s/api/bookmarklet?url='+encodeURIComponent(window.location.href)+'&title='+encodeURIComponent(document.title)+'&desc='+encodeURIComponent(metaDescription),'save-bookmark','width=500,height=300');
})());`
