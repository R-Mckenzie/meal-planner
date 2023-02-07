var deleteButtons: NodeListOf<HTMLElement> = document.querySelectorAll<HTMLElement>(".delete-button");
deleteButtons.forEach((b) => {
	const recipesContainer = document.querySelector<HTMLElement>("#recipesList")
	if (b.parentNode) {
		const node: HTMLElement = b.parentNode as HTMLElement
		b.addEventListener('click', () => {
			let rID: number = node.dataset.recipeid ? +node.dataset.recipeid : -1;
			fetch("/recipes", {
				method: "DELETE", body: JSON.stringify({ recipeID: +rID, csrf: recipesContainer?.dataset.csrf })
			})
			node.remove()
		})
	}
});

var recipeNode = document.querySelectorAll<HTMLElement>(".recipe-item");
recipeNode.forEach((n) => {
	const recipeID = n.dataset.recipeid;
	n.addEventListener('click', () => {
		location.assign("/recipes/" + recipeID)
	})
})
