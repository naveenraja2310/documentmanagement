<template>
    <div>
        <ul>
        <li class="left-item">
            <h3 class="app-url">Functionize</h3>
        </li>
        </ul>
    </div>  
</template>

<script>
</script>

<style scoped>
    ul {
    list-style-type: none;
    margin: 0;
    padding: 0;
    overflow: hidden;
    background-color: #333;
    display: flex; /* Use flexbox for layout */
    justify-content: space-between; /* Distribute items evenly across the width */
    }

    li {
    flex: 1; /* Equal width for both items */
    text-align: center; /* Center align text */
    }

    .app-url {
    color: #fff;
    text-transform: uppercase;
    border: 3px solid rgba(0, 178, 255, 0);
    border-radius: 20px;
    padding: 16px;
    font-family: Rubik, sans-serif;
    line-height: 0.9em; /* Adjust line height */
    text-decoration: none;
    transition: all 0.4s cubic-bezier(0.165, 0.84, 0.44, 1);
    display: block;
    }

    .app-url span {
        font-size: 16px;
        margin-bottom: 2px;

    }
    /* Specific styles for left and center items */
    .left-item {
    text-align: left; /* Align left */
    }

    .center-item {
    text-align: center; /* Align center */
    }
</style>
