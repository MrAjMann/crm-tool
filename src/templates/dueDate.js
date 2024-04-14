


document.addEventListener("DOMContentLoaded", function() {
    var dueDateInput = document.getElementById("invoice-DueDate");
    if (dueDateInput) {
        var date = new Date();
        date.setDate(date.getDate() + 30); // Add 30 days
        dueDateInput.value = date.toISOString().split('T')[0]; // Set the default date in YYYY-MM-DD format
    } else {
        console.error("DueDate input field not found.");
    }
});