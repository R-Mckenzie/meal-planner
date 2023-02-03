var recipes = document.querySelectorAll(".recipe-item");
var dropzones = document.querySelectorAll('.dropzone');
var dragged = null;
var addNodeListeners = function (node) {
    node.addEventListener('dragstart', function () {
        dragged = node;
        setTimeout(function () {
            if (dragged.classList.contains("meal-item")) {
                dragged.style.display = "none";
            }
        }, 0);
    });
    node.addEventListener('dragend', function () {
        setTimeout(function () {
            dragged.style.display = "flex";
            dragged = null;
        }, 0);
    });
};
recipes.forEach(function (r) {
    addNodeListeners(r);
});
var deleteButtons = document.querySelectorAll(".delete-button");
var addDeleteListener = function (b) {
    var node = b.parentNode;
    b.addEventListener('click', function () {
        if (node.classList.contains("meal-item")) {
            node.remove();
        }
        else {
            console.log("delete clicked");
        }
    });
};
deleteButtons.forEach(function (b) {
    addDeleteListener(b);
});
dropzones.forEach(function (z) {
    z.addEventListener('dragover', function (e) {
        e.preventDefault();
    });
    z.addEventListener('dragleave', function (e) {
        e.preventDefault();
    });
    z.addEventListener('drop', function (e) {
        if (dragged.classList.contains("meal-item")) {
            // If we are dragging from another date, move the node.
            z.append(dragged);
        }
        else {
            // If we are dragging from the recipe list, copy it.
            var nodeCopy = dragged.cloneNode(true);
            nodeCopy.id = Date.now();
            nodeCopy.classList.add("meal-item");
            addNodeListeners(nodeCopy);
            addDeleteListener(nodeCopy.querySelector('.delete-button'));
            z.append(nodeCopy);
        }
    });
});
// DATES =========================
function setToMonday(date) {
    var day = date.getDay() || 7;
    if (day !== 1)
        date.setHours(-24 * (day - 1));
    return date;
}
var date = document.getElementById("date");
var monday = setToMonday(new Date);
date.textContent = monday.toDateString();
// SAVING MEAL PLAN =================
// Returns the list of meals with their id and day of the week
var getMeals = function () {
    // {recipeID: x, date: xyz}
    var meals = [];
    dropzones.forEach(function (z) {
        z.childNodes.forEach(function (n) {
            if (!n.dataset)
                return;
            var date = new Date(monday);
            date.setDate(date.getDate() + parseInt(z.id));
            console.log(date);
            meals.push({ recipeID: +n.dataset.recipeid, date: date });
        });
    });
    return meals;
};
var saveButton = document.querySelector(".save-button");
saveButton.addEventListener("click", function () {
    var container = document.querySelector(".dashboard-container");
    fetch("/dashboard", {
        method: "POST", body: JSON.stringify({ meals: getMeals(), csrf: container.dataset.csrf })
    });
});
