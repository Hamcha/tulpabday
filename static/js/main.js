$(document).ready(function(){
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
		var empty = function (a) { return a.value.replace(/\./ig,"").length === 0; };
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
			localStorage.host = form.host.value;
			localStorage.secret = form.secret.value;
			location.reload();
		}).fail(function(xhr){
			$("#addTulpaPrompt").html("Something went wrong, if it persists, contact <a href=\"https://twitter.com/hamcha\">@hamcha</a>.<br />The error was: "+xhr.responseText);
		});
	});
	var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
	var suf = function(a) { return (a === 1 ? "st" : a === 2 ? "nd" : a === 3 ? "rd" : "th"); };
	var last = function(a) { return a[a.length - 1]; };
	$(".date").each(function(index,item){
		var date = new Date(item.innerHTML * 1000);
		item.innerHTML = "Born <span>" + months[date.getMonth()] + " " + date.getDate() + suf(last(date.getDate().toString())) + ", " + date.getFullYear() + "</span>";
	});
	$(".age").each(function(index,item){
		var date = new Date(item.innerHTML * 1000);
		item.innerHTML = "Age <span>"+age(date,false)+"</span>";
	});
	$(".until").each(function(index,item){
		var date = new Date(item.innerHTML * 1000);
		item.innerHTML = "Next birthday in "+age(date,true);
	});
});

function age(bday, offset) {
	var newDate = bday;
	var today = new Date();
	var difference;
	if (offset)
	{
		today.setYear(today.getFullYear()+1);
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
	if (years == 1) strout += years + " year, <br/>";
	else if (years > 0) strout += years + " years, ";
	if (months == 1) strout += months + " month, ";
	else if (months > 0) strout += months + " months, ";
	if (days == 1) strout += days + " day";
	else strout += days + " days";
	return strout;
}