const cardSearch = document.querySelector("input.search-bar")
const searchResults = document.querySelector("div.search-results")
const deckGrid = document.querySelector("div.display-deck")
const save = document.querySelector("button.save")
cardSearch.addEventListener('keypress', displaySearch)
deckGrid.addEventListener('click', deckClick)
save.addEventListener('click', saveDeck)


async function displaySearch(e){
    if(e.key === 'Enter'){
        searchResults.innerHTML = '';
        commander = cardSearch.value 
        try{
            const queryCommander = await getCardsByName(commander)

            if (queryCommander === null){
                window.alert("No results found")
                return
            }

            queryCommander.forEach(item => {
                //Card Holder DIV   
                const card = document.createElement("div")
                card.classList.add("card")
                card.dataset.index = item.id

                //The Card Image
                const image = document.createElement("img")
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
        }catch(error){
            console.log(error)
            window.alert("Network Error")
        }

    }
}



async function getCardsByName(name){
    try{
        const resp = await fetch(`https://localhost:4000/cards/search/${name}`, {
            method: 'GET', 
            mode: 'cors', 
        })
        return await resp.json()
    }catch(error){
        throw new Error("fetch failed")
    }
    
}

function addToDeck(e){

    //create the card holder
    const card = document.createElement("div")
    card.setAttribute("data-index", this.parentNode.getAttribute("data-index"))
    card.classList.add("card")

    const del = document.createElement("input")
    del.classList.add("delete-card") 
    del.type = "image"
    del.src = "/static/img/remove.png"

    const commander = document.createElement("input")
    commander.classList.add("make-commander")
    commander.type = "image"
    commander.src = "/static/img/commander.png"

    img = this.previousElementSibling
    newimg = document.createElement('img')
    newimg.src = img.src
    newimg.setAttribute("alt", img.getAttribute("alt"))
    newimg.width = 146;
    newimg.height = 204;

    card.append(del)
    card.append(commander)
    card.append(newimg)
    deckGrid.appendChild(card)
   
}

function deckClick(e){
    if (e.target.className === "delete-card"){
        e.target.parentNode.remove()
    } else if (e.target.className === "make-commander"){
        const previousCom = document.querySelector("div.commander")
        if (previousCom != null){
            previousCom.classList.remove("commander")
        }else{
            console.log("not found")
        }
        
        const card = e.target.parentNode;
        card.classList.add("commander")
    }
    
}

async function saveDeck(e){
    const deckGrid = document.querySelector('.display-deck')
    const items = deckGrid.querySelectorAll('div.card:not(.commander)')

    
    const itemIndexes = [...items].map(item => (item.getAttribute("data-index")))
    const itemIndexsNumeric = itemIndexes.map(Number)

    const com = document.querySelector(".commander")
    let comIndex = 0
    if (com != null){
        let temp = com.getAttribute("data-index")
        comIndex = parseInt(temp)
        if (isNaN(comIndex)){
            comIndex = 0
        }
    }

    const deckInfo = {
        'commander': comIndex, 
        'decklist': itemIndexsNumeric, 
    }; 

    try{
        const response = await fetch("https://localhost:4000/cards/save",{
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(deckInfo), 
        });
        if (!response.ok){
            console.log("Not Cool Man!")
        }
    }catch(error){
        console.log("Hello")
    }

}

