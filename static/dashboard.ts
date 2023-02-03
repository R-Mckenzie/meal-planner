let recipes: NodeListOf<HTMLElement> = document.querySelectorAll<HTMLElement>(".recipe-item");
let dropzones: NodeListOf<HTMLElement> = document.querySelectorAll<HTMLElement>('.dropzone');
let dragged: HTMLElement | null;

const addNodeListeners = (node: HTMLElement) => {
	node.addEventListener('dragstart', function() {
		dragged = node;
		setTimeout(function() {
			if (dragged?.classList.contains("meal-item")) {
				dragged.style.display = "none";
			}
		}, 0)
	});

	node.addEventListener('dragend', function() {
		setTimeout(function() {
			if (dragged) {
				dragged.style.display = "flex";
				dragged = null;
			}
		}, 0);
	});
}

recipes.forEach((r) => {
	addNodeListeners(r);
});

const addDeleteListener = (b: HTMLElement) => {
	if (b.parentNode) {
		const node: HTMLElement = b.parentNode as HTMLElement
		b.addEventListener('click', () => {
			if (node.classList.contains("meal-item")) {
				node.remove()
			} else {
				console.log("delete clicked")
			}
		})
	}
}

let deleteButtons: NodeListOf<HTMLElement> = document.querySelectorAll<HTMLElement>(".delete-button");
deleteButtons.forEach((b) => {
	addDeleteListener(b)
});

dropzones.forEach((z) => {
	z.addEventListener('dragover', function(e) {
		e.preventDefault();
	});
	z.addEventListener('dragleave', function(e) {
		e.preventDefault();
	});
	z.addEventListener('drop', function(e) {
		if (!dragged) return;
		if (dragged?.classList.contains("meal-item")) {
			// If we are dragging from another date, move the node.
			z.append(dragged)
		} else {
			// If we are dragging from the recipe list, copy it.
			let nodeCopy: HTMLElement = dragged.cloneNode(true) as HTMLElement;
			nodeCopy.id = Date.now().toString();
			nodeCopy.classList.add("meal-item")
			addNodeListeners(nodeCopy);
			const deleteBtn = nodeCopy.querySelector<HTMLElement>('.delete-button')
			if (deleteBtn) addDeleteListener(deleteBtn)
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

let date = document.getElementById("date")
let monday = setToMonday(new Date);
if (date) date.textContent = monday.toDateString();

// SAVING MEAL PLAN =================

// Returns the list of meals with their id and day of the week
const getMeals = () => {
	// {recipeID: x, date: xyz}
	let meals = []
	dropzones.forEach((z) => {
		z.childNodes.forEach((n) => {
			if (!n.dataset) return;
			let date = new Date(monday);
			date.setDate(date.getDate() + parseInt(z.id));
			console.log(date)
			meals.push({ recipeID: +n.dataset.recipeid, date: date })
		})
	})
	return meals;
}

const saveButton = document.querySelector(".save-button")
saveButton.addEventListener("click", () => {
	const container = document.querySelector(".dashboard-container")
	fetch("/dashboard", {
		method: "POST", body: JSON.stringify({ meals: getMeals(), csrf: container.dataset.csrf })
	})
})
