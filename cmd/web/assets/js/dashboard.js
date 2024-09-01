htmx.onLoad(function(content) {
	var recipes = content.querySelectorAll(".recipes");
	for (var i = 0; i < recipes.length; i++) {
		var recipe = recipes[i]
		recipesInstance = new Sortable(recipe, {
			sort: false,
			group: {
				name: 'meals',
				pull: 'clone',
				put: false
			},
			animation: 100,
		});

	}

	var mealLists = content.querySelectorAll(".sortable");
	for (var i = 0; i < mealLists.length; i++) {
		var sortable = mealLists[i];
		var sortableInstance = new Sortable(sortable, {
			group: 'meals',
			animation: 100,

			// Make the `.htmx-indicator` unsortable
			filter: ".htmx-indicator",
			onMove: function(evt) {
				return evt.related.className.indexOf('htmx-indicator') === -1;
			},
		});
	}
})
