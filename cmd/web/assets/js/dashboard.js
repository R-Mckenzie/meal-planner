htmx.onLoad(function(content) {
	var sortables = content.querySelectorAll(".sortable");
	for (var i = 0; i < sortables.length; i++) {
		var sortable = sortables[i];
		var sortableInstance = new Sortable(sortable, {
			group: 'meals',
			animation: 150,
			ghostClass: "text-foreground/50",

			// Make the `.htmx-indicator` unsortable
			filter: ".htmx-indicator",
			onMove: function(evt) {
				return evt.related.className.indexOf('htmx-indicator') === -1;
			},

			// Disable sorting on the `end` event
			onEnd: function(evt) {
				this.option("disabled", true);
			}
		});

		// Re-enable sorting on the `htmx:afterSwap` event
		sortable.addEventListener("htmx:afterSwap", function() {
			sortableInstance.option("disabled", false);
		});
	}

	var recipes = content.querySelectorAll(".recipes");
	for (var i = 0; i < recipes.length; i++) {
		var recipe = recipes[i]
		var recipeInstance = new Sortable(recipe, {
			sort: false,
			group: {
				name: 'meals',
				pull: 'clone',
				put: false
			},
			animation: 150,
			ghostClass: "text-foreground/50",

			// Make the `.htmx-indicator` unsortable
			filter: ".htmx-indicator",
			onMove: function(evt) {
				return evt.related.className.indexOf('htmx-indicator') === -1;
			},
		});

		// Re-enable sorting on the `htmx:afterSwap` event
		sortable.addEventListener("htmx:afterSwap", function() {
			recipeInstance.option("disabled", false);
		});
	}
})
