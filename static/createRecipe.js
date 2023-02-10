var redirectLocField = document.getElementById("return_url");
var urlParams = new URLSearchParams(location.search);
redirectLocField.value = urlParams.get("returnto") || "dashboard";
