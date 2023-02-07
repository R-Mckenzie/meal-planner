let recipes = document.querySelectorAll(".recipe-item") as NodeListOf<HTMLElement>;
let dropzones = document.querySelectorAll('.dropzone') as NodeListOf<HTMLElement>;
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

var addDeleteListener = (b: HTMLElement) => {
	const container = document.querySelector<HTMLElement>(".dashboard-container")
	if (b.parentNode) {
		const node: HTMLElement = b.parentNode as HTMLElement
		b.addEventListener('click', () => {
			if (node.classList.contains("meal-item")) {
				node.remove()
			} else {
				let rID: number = node.dataset.recipeid ? +node.dataset.recipeid : -1;
				fetch("/recipes", {
					method: "DELETE", body: JSON.stringify({ recipeID: +rID, csrf: container?.dataset.csrf })
				})
				node.remove()
			}
		})
	}
}

var deleteButtons: NodeListOf<HTMLElement> = document.querySelectorAll<HTMLElement>(".delete-button");
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
function setToMonday(date: Date) {
	var day = date.getDay() || 7;
	if (day !== 1)
		date.setHours(-24 * (day - 1));
	return date;
}

let dateElement = document.getElementById("date")
let monday = setToMonday(new Date(Date.parse(dateElement?.textContent || "")));

const changeDate = () => {
	// const queryDate = monday.getDate() + "-" + (monday.getMonth() + 1). + "-" + monday.getFullYear()
	const queryDate = monday.toISOString().slice(0, 10)
	window.location.replace("/dashboard?date=" + queryDate)
}

const prevBtn = document.getElementById("prev-week") as HTMLElement
const nextBtn = document.getElementById("next-week") as HTMLElement

prevBtn.addEventListener("click", function() {
	monday.setDate(monday.getDate() - 7)
	changeDate()
})

nextBtn.addEventListener("click", function() {
	monday.setDate(monday.getDate() + 7)
	changeDate()
})


// SAVING MEAL PLAN =================
type Meal = {
	recipeID: number,
	date: Date,
}

const getMeals = () => {
	let meals: Meal[] = []
	dropzones.forEach((z: HTMLElement) => {
		z.childNodes.forEach((node) => {
			let n = node as HTMLElement
			if (!n.dataset)
				return;
			let date = new Date(+monday);
			date.setDate(date.getDate() + parseInt(z.id));

			let rID: number = n.dataset.recipeid ? +n.dataset.recipeid : -1;
			meals.push({ recipeID: rID, date: date })
		})
	})
	return meals;
}

const saveButton = document.querySelector(".save-button") as HTMLElement
saveButton.addEventListener("click", () => {
	const container = document.querySelector<HTMLElement>(".dashboard-container")
	fetch("/dashboard", {
		method: "POST", body: JSON.stringify({ weekStart: monday, meals: getMeals(), csrf: container?.dataset.csrf })
	})
})
