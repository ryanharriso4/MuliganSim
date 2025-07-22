const cardSearch = document.querySelector("input.search-bar")
const searchResults = document.querySelector("div.search-results")
cardSearch.addEventListener('keypress', displaySearch)


async function displaySearch(e){
    if(e.key === 'Enter'){
        searchResults.innerHTML = '';
        commander = cardSearch.value 
        const queryCommander = await getCardsByName(commander)
        queryCommander.forEach(item => {
            //Card Holder DIV   
            const card = document.createElement("div")
            card.className = "card"

            //The Card Image
            const image = document.createElement("img")
            image.dataset.index = item.id
            image.classList.add("card-image")  
            image.src = item.imageuri
            image.alt = item.name

            //Button for adding card to deck
            const deckAdd = document.createElement("input")
            deckAdd.classList.add("add-card") 
            deckAdd.type = "image"
            deckAdd.src = "/static/img/add.png"
            deckAdd.addEventListener('click', addToDeck)

            card.append(image)
            card.append(deckAdd)

            searchResults.append(card)
        })
    }
}



async function getCardsByName(name){
    const resp = await fetch(`https://localhost:4000/cards/search/${name}`, {
        method: 'GET', 
        mode: 'cors', 
    })
    return await resp.json()
}

function addToDeck(e){
    img = this.previousElementSibling
    newimg = document.createElement('img')
    newimg.src = img.src
    newimg.setAttribute("data-index", img.getAttribute("data-index"))
    newimg.setAttribute("alt", img.getAttribute("alt"))
    buildgrid = document.querySelector("div.display-deck")
    buildgrid.appendChild(newimg)
   
}




