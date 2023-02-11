
// Get the nosurf csrf token
const csrfToken = document.querySelector<HTMLElement>(".dashboard-container")?.dataset.csrf

/* DRAG AND DROP
 * 	- Adds drag and drop event listeners for nodes and dropzones
 * 	- Copying nodes from recipe list, moving nodes already on calendar
 */
let dragged: HTMLElement | null; // Dragged is the currently held drag and drop node
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

let recipes = document.querySelectorAll(".recipe-item") as NodeListOf<HTMLElement>; // Includes the recipe list and meals on calendar
recipes.forEach((r) => {
	addNodeListeners(r);
});

// If the node is in the recipe list, delete from database.
// If node is in the calendar, delete it from the page but not the database
var addDeleteListener = (b: HTMLElement) => {
	if (b.parentNode) {
		const node: HTMLElement = b.parentNode as HTMLElement
		b.addEventListener('click', () => {
			if (node.classList.contains("meal-item")) {
				node.remove()
			} else {
				let rID: number = node.dataset.recipeid ? +node.dataset.recipeid : -1;
				fetch("/recipes", {
					method: "DELETE", body: JSON.stringify({ recipeID: +rID, csrf: csrfToken })
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

let dropzones = document.querySelectorAll('.dropzone') as NodeListOf<HTMLElement>;
dropzones.forEach((z) => {
	z.addEventListener('dragover', function(e) {
		e.preventDefault();
	});
	z.addEventListener('dragleave', function(e) {
		e.preventDefault();
	});
	// If we are dragging from another dropzone, move the node.
	// If we are dragging from the recipe list, copy it.
	z.addEventListener('drop', function(e) {
		if (!dragged) return;
		if (dragged?.classList.contains("meal-item")) {
			z.append(dragged)
		} else {
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

/* DATES 
 * 	- handles getting the date fromt the date elemtent,
 * 	  which comes from the backend, and getting the first day of that week (Monday)
 * 	- Moving forward and back a week with the prev and next buttons. 
 * 	  This sends a request to the backend to reload the page with the new week's meals
 */
const setToMonday = (date: Date) => {
	var day = date.getDay() || 7;
	if (day !== 1)
		date.setHours(-24 * (day - 1));
	return date;
}

let dateElement = document.getElementById("date")
let monday = setToMonday(new Date(Date.parse(dateElement?.textContent || "")));

const changeDate = () => {
	const queryDate = monday.toISOString().slice(0, 10)
	// Save meals?
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


/* SAVING MEALS 
 * 	- Scans the dropzones and gets node recipeIDs and their respective dates
 * 	- Sends a post request with the meal data and CSRF token for the backend to save
 */
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
	fetch("/dashboard", {
		method: "POST", body: JSON.stringify({ weekStart: monday, meals: getMeals(), csrf: csrfToken })
	}).then(() => location.reload())
})
