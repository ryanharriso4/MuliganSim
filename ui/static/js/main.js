
const addbuttons = document.querySelectorAll("input.add-card")
addbuttons.forEach(function(currentBtn){
    currentBtn.addEventListener('click', addToDeck)
})

function addToDeck(e){
    img = this.previousElementSibling
    newimg = document.createElement('img')
    newimg.src = img.src
    newimg.setAttribute("data-index", img.getAttribute("data-index"))
    newimg.setAttribute("alt", img.getAttribute("alt"))
    buildgrid = document.querySelector("div.display-deck")
    buildgrid.appendChild(newimg)
   
}






