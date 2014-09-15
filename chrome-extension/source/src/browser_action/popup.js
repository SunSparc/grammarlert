var selectedText;

function updateFormAction(event) {
	chrome.tabs.query({currentWindow: true, active: true}, function(tabs)
	{
		var url = tabs[0].url;
		var suggestion = document.getElementById("suggestion").value;
  		var formUrl = "http://grammarlert.com/app/report";
  		var params = '{'+
  						'"url":"'+url+'",'+
  						'"original":"'+selectedText+'",'+
  						'"suggest":"'+suggestion+'"'+
  					'}'; // TODO:add note
		var xhr = new XMLHttpRequest();
		xhr.open("POST", formUrl, true);
		xhr.setRequestHeader("Content-type", "application/json; charset=utf-8")

		xhr.onreadystatechange = function()
		{
			if (xhr.readyState == 4)
				window.close();
		}

		xhr.send(params);
	});
	
	event.preventDefault();
    return false;
}


document.getElementById("form").addEventListener('submit', updateFormAction);


chrome.extension.onMessage.addListener(function(message, sender, sendResponse)
{
	selectedText = message.selection;
});

chrome.tabs.executeScript(null, {code: "chrome.extension.sendMessage({selection: window.getSelection().toString() });"});	