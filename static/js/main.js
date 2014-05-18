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
			name: form.name.value,
			host: form.host.value,
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
});
