const redirectLocField = document.getElementById("return_url") as HTMLInputElement
const urlParams = new URLSearchParams(location.search)
redirectLocField.value = urlParams.get("returnto") || "dashboard"
