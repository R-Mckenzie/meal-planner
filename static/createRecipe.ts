// Redirect to where the user came from on submit based on url parameter
const redirectLocField = document.getElementById("return_url") as HTMLInputElement
const urlParams = new URLSearchParams(location.search)
redirectLocField.value = urlParams.get("returnto") || "dashboard"
