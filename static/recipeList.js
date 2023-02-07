var deleteButtons = document.querySelectorAll(".delete-button");
deleteButtons.forEach(function (b) {
    var recipesContainer = document.querySelector("#recipesList");
    if (b.parentNode) {
        var node_1 = b.parentNode;
        b.addEventListener('click', function () {
            var rID = node_1.dataset.recipeid ? +node_1.dataset.recipeid : -1;
            fetch("/recipes", {
                method: "DELETE", body: JSON.stringify({ recipeID: +rID, csrf: recipesContainer === null || recipesContainer === void 0 ? void 0 : recipesContainer.dataset.csrf })
            });
            node_1.remove();
        });
    }
});
var recipeNode = document.querySelectorAll(".recipe-item");
recipeNode.forEach(function (n) {
    var recipeID = n.dataset.recipeid;
    n.addEventListener('click', function () {
        location.assign("/recipes/" + recipeID);
    });
});
