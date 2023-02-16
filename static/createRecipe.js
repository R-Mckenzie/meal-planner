// Redirect to where the user came from on submit based on url parameter
var redirectLocField = document.getElementById("return_url");
var urlParams = new URLSearchParams(location.search);
redirectLocField.value = urlParams.get("returnto") || "dashboard";
