var endpointUrl = "https://xxx"

function onEditInstalled(e){
  
  var sheet = e.range.getSheet()
  var sheetId = sheet.getParent().getId()
  var sheetName = sheet.getName()
    
  var url = endpointUrl+ "/check/"+sheetId+"/"+sheetName
  
  try {
    var resp = UrlFetchApp.fetch(url);
    var msg = resp.getContentText()
    if (msg != "") {
      Browser.msgBox(msg)
    }
  } catch(ex) {
    Browser.msgBox("server error\\nplease contact the system administrator")
  }
}