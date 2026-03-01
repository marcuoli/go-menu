package menu

import "fmt"

// linkClasses returns the Tailwind CSS classes for a menu link item.
// Active items get a highlighted background; inactive items get hover effects.
func linkClasses(active bool) string {
	base := "flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors"
	if active {
		return base + " bg-blue-50 text-blue-700 font-medium dark:bg-blue-900/30 dark:text-blue-300"
	}
	return base + " text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700/50"
}

// openState returns the Alpine.js x-data expression for a collapsible section.
// If the section has an active child, it starts open; otherwise collapsed.
func openState(item MenuItem, config MenuConfig) string {
	open := HasActiveChild(item)
	return fmt.Sprintf("{ open: %t }", open)
}
