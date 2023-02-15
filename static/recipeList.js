// Listener to delete recipe from database when clicked
var deleteButtons = document.querySelectorAll(".delete-button");
deleteButtons.forEach(function (b) {
    var _a;
    var csrfToken = (_a = document.querySelector("#recipesList")) === null || _a === void 0 ? void 0 : _a.dataset.csrf;
    if (b.parentNode) {
        var node_1 = b.parentNode;
        b.addEventListener('click', function () {
            var rID = node_1.dataset.recipeid ? +node_1.dataset.recipeid : -1;
            fetch("/recipes", {
                method: "DELETE", body: JSON.stringify({ recipeID: +rID, csrf: csrfToken })
            });
            node_1.remove();
        });
    }
});
// Listener to redirect to the recipes edit page
var recipeNode = document.querySelectorAll(".recipe-item");
recipeNode.forEach(function (n) {
    var recipeID = n.dataset.recipeid;
    n.addEventListener('click', function () {
        location.assign("/recipes/" + recipeID);
    });
});
