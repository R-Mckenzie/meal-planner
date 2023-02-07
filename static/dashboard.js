var recipes = document.querySelectorAll(".recipe-item");
var dropzones = document.querySelectorAll('.dropzone');
var dragged;
var addNodeListeners = function (node) {
    node.addEventListener('dragstart', function () {
        dragged = node;
        setTimeout(function () {
            if (dragged === null || dragged === void 0 ? void 0 : dragged.classList.contains("meal-item")) {
                dragged.style.display = "none";
            }
        }, 0);
    });
    node.addEventListener('dragend', function () {
        setTimeout(function () {
            if (dragged) {
                dragged.style.display = "flex";
                dragged = null;
            }
        }, 0);
    });
};
recipes.forEach(function (r) {
    addNodeListeners(r);
});
var addDeleteListener = function (b) {
    var container = document.querySelector(".dashboard-container");
    if (b.parentNode) {
        var node_1 = b.parentNode;
        b.addEventListener('click', function () {
            if (node_1.classList.contains("meal-item")) {
                node_1.remove();
            }
            else {
                var rID = node_1.dataset.recipeid ? +node_1.dataset.recipeid : -1;
                fetch("/recipes", {
                    method: "DELETE", body: JSON.stringify({ recipeID: +rID, csrf: container === null || container === void 0 ? void 0 : container.dataset.csrf })
                });
                node_1.remove();
            }
        });
    }
};
var deleteButtons = document.querySelectorAll(".delete-button");
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
        if (!dragged)
            return;
        if (dragged === null || dragged === void 0 ? void 0 : dragged.classList.contains("meal-item")) {
            // If we are dragging from another date, move the node.
            z.append(dragged);
        }
        else {
            // If we are dragging from the recipe list, copy it.
            var nodeCopy = dragged.cloneNode(true);
            nodeCopy.id = Date.now().toString();
            nodeCopy.classList.add("meal-item");
            addNodeListeners(nodeCopy);
            var deleteBtn = nodeCopy.querySelector('.delete-button');
            if (deleteBtn)
                addDeleteListener(deleteBtn);
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
var dateElement = document.getElementById("date");
var monday = setToMonday(new Date(Date.parse((dateElement === null || dateElement === void 0 ? void 0 : dateElement.textContent) || "")));
var changeDate = function () {
    // const queryDate = monday.getDate() + "-" + (monday.getMonth() + 1). + "-" + monday.getFullYear()
    var queryDate = monday.toISOString().slice(0, 10);
    window.location.replace("/dashboard?date=" + queryDate);
};
var prevBtn = document.getElementById("prev-week");
var nextBtn = document.getElementById("next-week");
prevBtn.addEventListener("click", function () {
    monday.setDate(monday.getDate() - 7);
    changeDate();
});
nextBtn.addEventListener("click", function () {
    monday.setDate(monday.getDate() + 7);
    changeDate();
});
var getMeals = function () {
    var meals = [];
    dropzones.forEach(function (z) {
        z.childNodes.forEach(function (node) {
            var n = node;
            if (!n.dataset)
                return;
            var date = new Date(+monday);
            date.setDate(date.getDate() + parseInt(z.id));
            var rID = n.dataset.recipeid ? +n.dataset.recipeid : -1;
            meals.push({ recipeID: rID, date: date });
        });
    });
    return meals;
};
var saveButton = document.querySelector(".save-button");
saveButton.addEventListener("click", function () {
    var container = document.querySelector(".dashboard-container");
    fetch("/dashboard", {
        method: "POST", body: JSON.stringify({ weekStart: monday, meals: getMeals(), csrf: container === null || container === void 0 ? void 0 : container.dataset.csrf })
    });
});
