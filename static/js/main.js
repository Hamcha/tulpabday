$(document).ready(function(){
	// Put localStorage data if saved
	var addform = document.getElementById("addTulpaForm");
	if (typeof localStorage.host !== "undefined")
		addform.host.value = localStorage.host;
	if (typeof localStorage.secret !== "undefined")
		addform.secret.value = localStorage.secret;

	// Show "Add tulpa" form
	$("#addTulpaLink").click(function(e){
		$("#addTulpaPrompt").slideToggle();
	});

	// Add tulpa via XHR
	$("#addTulpaBtn").click(function(e){
		var form = addform;
		// Get field list
		var fields = [form.name, form.host, form.date, form.secret];
		// Clear previous red borders
		for (var i = 0; i < fields.length; i++)
			fields[i].className = fields[i].className.replace(" red","");
		// Check for empty fields
		fields = fields.filter(empty);
		if (fields.length > 0) {
			for (i = 0; i < fields.length; i++)
				fields[i].className += " red";
			return;
		}
		$.post("/add",{
			name: form.name.value.replace(/\./ig,""),
			host: form.host.value.replace(/\./ig,""),
			date: Date.parse(form.date.value)/1000|0,
			secret: sha256_digest(form.secret.value)
		}).done(function(data){
			// Save to localStorage
			localStorage.host = form.host.value;
			localStorage.secret = form.secret.value;
			// Reload page
			location.reload();
		}).fail(function(xhr){
			$("#addTulpaPrompt").html("Something went wrong, if it persists, contact <a href=\"https://twitter.com/hamcha\">@hamcha</a>.<br />The error was: "+xhr.responseText);
		});
	});

	// Generate tulpa list
	var out = "";
	var imminent = false;
	var immiout = "";
	var bday = false;
	var bout = "";
	var i = 0;
	var now = new Date();
	for (i = 0; i < data.length; i++) {
		var obj  = data[i];
		var date = new Date(obj.Born * 1000);
		var id   = (obj.Name+"_"+obj.Host).replace(" ", "_");
		tulpas[id] = obj;
		out += "<div class=\"tulpa\" id=\""+id+"\">"+render(obj)+"</div>";
		// Check for imminent birthday (under 30 days)
		var nextdate = new Date(date);
		nextdate.setYear(now.getFullYear());
		var daysdiff = (nextdate - now)/1000/60/60/24;
		var newage = 0;
		if (daysdiff > 0 && daysdiff < 30) {
			imminent = true;
			newage = Math.ceil((now-date) / (1000*60*60*24*365));
			immiout += "<li class=\"tulpali\"><span class=\"name\">"+obj.Name+"</span> "+
			"<span class=\"host\">"+obj.Host+"</span> will turn "+
			"<span class=\"age\">"+newage+"</span> in "+
			"<span class=\"date\">"+Math.ceil(daysdiff)+"</span> days"+
			"</li>";
		}
		// Check for current birthdays
		daysdiff = (nextdate - now)/1000/60/60/24+1;
		if (daysdiff < 1 && daysdiff > 0) {
			bday = true;
			newage = Math.floor((now-date) / (1000*60*60*24*365));
			bout += "<li class=\"tulpali\"><span class=\"name\">"+obj.Name+"</span> "+
			"<span class=\"host\">"+obj.Host+"</span> has turned "+
			"<span class=\"age\">"+newage+"</span> today!"+
			"</li>";
		}
	}

	$("#tulpas").html(out);

	if (!imminent) {
		$(".hideim").remove();
	} else {
		$("#imminent").html(immiout);
		$("#imminent .tulpali").sortElements(imsort);
	}

	if (!bday) {
		$(".hidebday").remove();
	} else {
		$("#bday").html(bout);
	}

	// Order tulpas by bday
	$(".tulpa").sortElements(tbsort);
});

var tulpas = {};

// Util snippets
var empty = function (a) { return a.value.replace(/\./ig,"").length === 0; };
var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
var suf = function(a) {
	switch (a) {
		case "21": case "1": return "st";
		case "22": case "2": return "nd";
		case "23": case "3": return "rd";
		default:             return "th";
	}
};
var last = function(a) { return a[a.length - 1]; };

var render = function(obj) {
	var date = new Date(obj.Born * 1000);
	var suffix = suf(date.getDate().toString());
	var born = "Born <span>" + months[date.getMonth()] + " " + date.getDate() + suffix + ", " + date.getFullYear() + "</span>";
	var tage = "Age <span>"+age(date,false)+"</span>";
	var next = "Next birthday in "+age(date,true);
	var id   = (obj.Name+"_"+obj.Host).replace(" ", "_");
	return "<table width=\"100%\"><tr><td>"+
			"<div class=\"name\">"+obj.Name+"</div>"+
			"<div class=\"host\">"+obj.Host+"</div>"+
			"</td><td style=\"width: 140px;\">"+
			"<div class=\"time\" style=\"display:none;\">"+obj.Born+"</div>"+
			"<div class=\"date\">"+born+"</div>"+
			"<div class=\"age\">"+tage+"</div></td></tr>"+
			"<tr><td colspan=\"2\" style=\"text-align: center;\">"+
			"<div class=\"until\">"+next+"</div>"+
			"<div class=\"edit\" onclick=\"editprompt('"+id+"')\">&nbsp</div><div class=\"delete\" onclick=\"deleteprompt('"+id+"')\">&nbsp</div>"+
			"</td></tr></table>";
};

var pad = function (x,n) {
	return ("0".repeat(n) + x.toString()).slice(-n);
}

var showEditDate = function (unixtime) {
	var date = new Date(unixtime*1000);
	return date.getUTCFullYear() + "-" + pad(date.getUTCMonth() + 1, 2) + "-" + pad(date.getUTCDate(), 2);
};

var editprompt = function (id) {
	var element = $("#"+id);
	element.addClass("prompt");
	var obj = tulpas[id];
	var prompt = "<form id=\"EDIT_"+id+"\" onsubmit=\"return false;\">"+
		"<input type=\"text\" class=\"short\" name=\"name\"\" placeholder=\"Tulpa's name\" value=\""+obj.Name+"\">&nbsp;"+
		"<input type=\"text\" class=\"short\" name=\"host\"\" placeholder=\"Host's name\" value=\""+obj.Host+"\"><br />"+
		"<input type=\"date\" class=\"short\" name=\"date\"\" placeholder=\"Birth Date\" value=\""+showEditDate(obj.Born)+"\">&nbsp;"+
		"<input type=\"password\" class=\"short\" name=\"secret\"\" placeholder=\"Secret code\"><br />"+
		"<button onclick=\"tryEdit('"+id+"','EDIT_"+id+"')\">Edit</button><button onclick=\"restore('"+id+"')\">Nevermind</button></form>";
	$("#"+id).html(prompt);
};

var deleteprompt = function (id) {
	var element = $("#"+id);
	element.addClass("prompt");
	var prompt = "Write the secret code to <br /><b style=\"color: #c07;\">delete</b> your tulpa<br />"+
		"<input type=\"password\" id=\"DELETE_"+id+"\"><br />"+
		"<button onclick=\"tryDelete('"+id+"','DELETE_"+id+"')\">Delete</button><button onclick=\"restore('"+id+"')\">Nevermind</button>";
	$("#"+id).html(prompt);
};

var restore = function (id) {
	$("#"+id).removeClass("prompt").html(render(tulpas[id]));
};

var tryDelete = function (id,codeid) {
	var code = $("#"+codeid).val();
	if (code.length < 1) return $("#"+codeid).addClass("red");
	codesha = sha256_digest(code);
	obj = tulpas[id];
	$.post("/delete",{ name: obj.Name, host: obj.Host, secret: codesha
	}).done(function (data){
		location.reload();
	}).fail(function (xhr){
		$("#"+id).html("<b>Failed</b> : "+xhr.responseText + "<br /><button onclick=\"restore('"+id+"')\">Nevermind</button>");
	});
};

var tryEdit = function (id,codeid) {
	var form = document.getElementById(codeid);
	// Get field list
	var fields = [form.name, form.host, form.date, form.secret];
	// Clear previous red borders
	for (var i = 0; i < fields.length; i++)
		fields[i].className = fields[i].className.replace(" red","");
	// Check for empty fields
	fields = fields.filter(empty);
	if (fields.length > 0) {
		for (i = 0; i < fields.length; i++)
			fields[i].className += " red";
		return;
	}
	obj = tulpas[id];
	// Send XHR (and hope for the best)
	$.post("/edit",{
			oldname: obj.Name,
			oldhost: obj.Host,
			name: form.name.value.replace(/\./ig,""),
			host: form.host.value.replace(/\./ig,""),
			date: Date.parse(form.date.value)/1000|0,
			secret: sha256_digest(form.secret.value)
	}).done(function (data){
		location.reload();
	}).fail(function (xhr){
		$("#"+id).html("<b>Failed</b> : "+xhr.responseText + "<br /><button onclick=\"restore('"+id+"')\">Nevermind</button>");
	});
};

var parseDate = function (date) {
	var d = new Date(date*1000);
	return [d.getUTCFullYear(),d.getUTCMonth()+1,d.getUTCDate()].join("/");
};

var tbsort = function (a,b) {
	var dateA = new Date($(a).find(".time")[0].innerHTML * 1000);
	var dateB = new Date($(b).find(".time")[0].innerHTML * 1000);
	var tod = new Date();
	tod.setYear(tod.getFullYear()+1);
	ad = tod-dateA; bd = tod-dateB;
	ay = Math.ceil(ad / (1000*60*60*24*365));
	by = Math.ceil(bd / (1000*60*60*24*365));
	ad -= ay*1000*60*60*24*365; bd -= by*1000*60*60*24*365;
	return bd - ad;
};

var imsort = function (a,b) {
	var dayA = parseInt($(a).find(".date")[0].innerHTML);
	var dayB = parseInt($(b).find(".date")[0].innerHTML);
	return dayA - dayB;
};

var age = function (bday, offset) {
	var newDate = bday;
	var today = new Date();
	var difference;
	if (offset)
	{
		today.setYear(today.getUTCFullYear()+1);
		difference = (newDate-today);
	} else {
		difference = (today-newDate);
	}
	var years = Math.floor(difference / (1000*60*60*24*365));
	difference -= years * (1000*60*60*24*365);
	var months = Math.floor(difference / (1000*60*60*24*30.4375));
	difference -= months * (1000*60*60*24*30.4375);
	var days = 0;
	if (offset)
		days = Math.ceil(difference / (1000*60*60*24));
	else
		days = Math.floor(difference / (1000*60*60*24));
	var strout = "";
	if (years == 1) strout += years + " year, <br />";
	else if (years > 0) strout += years + " years, <br />";
	if (months == 1) strout += months + " month, ";
	else if (months > 0) strout += months + " months, ";
	if (days == 1) strout += days + " day";
	else strout += days + " days";
	return strout;
};
