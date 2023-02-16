var _a;
// Get the nosurf csrf token
var csrfToken = (_a = document.querySelector(".dashboard-container")) === null || _a === void 0 ? void 0 : _a.dataset.csrf;
// GET REQUIRED HTML ELEMENTS
var dropzones = document.querySelectorAll('.dropzone');
var dateElement = document.getElementById("date");
var setToMonday = function (date) {
    var day = date.getDay() || 7;
    if (day !== 1)
        date.setHours(-24 * (day - 1));
    return date;
};
var monday = setToMonday(new Date(Date.parse((dateElement === null || dateElement === void 0 ? void 0 : dateElement.textContent) || "")));
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
var startingData = getMeals();
console.log(startingData);
var saveMeals = function (weekBeginning) {
    var newMeals = getMeals();
    console.debug(newMeals);
    if (startingData != newMeals) {
        console.log("saved");
        return fetch("/dashboard", { method: "POST", body: JSON.stringify({ weekStart: weekBeginning, meals: newMeals, csrf: csrfToken }) });
    }
    console.log("not saved");
};
var saveButton = document.querySelector(".save-button");
saveButton.addEventListener("click", function () {
    var _a;
    (_a = saveMeals(monday)) === null || _a === void 0 ? void 0 : _a.then(function () { return location.reload(); });
});
/* DRAG AND DROP
 * 	- Adds drag and drop event listeners for nodes and dropzones
 * 	- Copying nodes from recipe list, moving nodes already on calendar
 */
var dragged; // Dragged is the currently held drag and drop node
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
var recipes = document.querySelectorAll(".recipe-item"); // Includes the recipe list and meals on calendar
recipes.forEach(function (r) {
    addNodeListeners(r);
});
// If the node is in the recipe list, delete from database.
// If node is in the calendar, delete it from the page but not the database
var addDeleteListener = function (b) {
    if (b.parentNode) {
        var node_1 = b.parentNode;
        b.addEventListener('click', function () {
            if (node_1.classList.contains("meal-item")) {
                node_1.remove();
            }
            else {
                var rID = node_1.dataset.recipeid ? +node_1.dataset.recipeid : -1;
                fetch("/recipes/delete", {
                    method: "DELETE", body: JSON.stringify({ recipeID: +rID, csrf: csrfToken })
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
    // If we are dragging from another dropzone, move the node.
    // If we are dragging from the recipe list, copy it.
    z.addEventListener('drop', function (e) {
        if (!dragged)
            return;
        if (dragged === null || dragged === void 0 ? void 0 : dragged.classList.contains("meal-item")) {
            z.append(dragged);
        }
        else {
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
/* DATES
 * 	- handles getting the date fromt the date elemtent,
 * 	  which comes from the backend, and getting the first day of that week (Monday)
 * 	- Moving forward and back a week with the prev and next buttons.
 * 	  This sends a request to the backend to reload the page with the new week's meals
 */
var changeDate = function () {
    var queryDate = monday.toISOString().slice(0, 10);
    window.location.replace("/dashboard?date=" + queryDate);
};
var prevBtn = document.getElementById("prev-week");
var nextBtn = document.getElementById("next-week");
prevBtn.addEventListener("click", function () {
    saveMeals(monday);
    monday.setDate(monday.getDate() - 7);
    changeDate();
});
nextBtn.addEventListener("click", function () {
    saveMeals(monday);
    monday.setDate(monday.getDate() + 8);
    changeDate();
});
