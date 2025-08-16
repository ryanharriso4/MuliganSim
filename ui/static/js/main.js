const cardSearch = document.querySelector("input.search-bar")
const searchResults = document.querySelector("div.search-results")
const deckGrid = document.querySelector("div.display-deck")
const save = document.querySelector("button.save")

const del = document.querySelector("button.deleteDeck")

const deckAdd = []
const deckRemove = []

const lazyLoad = target => {
    const io = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            console.log("hello world")

            if (entry.isIntersecting){
                const img = entry.target
                const src = img.getAttribute('data-src')

                img.setAttribute('src', src)

                observer.disconnect()
            }
        })
    })
    io.observe(target)
}

if(deckGrid != null){
    const cards = deckGrid.querySelectorAll(".deck-card-image")
    cards.forEach(lazyLoad)
}


if (cardSearch != null) {
    cardSearch.addEventListener('keypress', displaySearch)
}

if (deckGrid != null){
    deckGrid.addEventListener('click', deckClick)
}

if (save != null){
    save.addEventListener('click', saveDeck)
}

 if (del != null){
    del.addEventListener('click', deleteClick)
 }



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
                image.setAttribute("data-src", item.imageuri)
                image.alt = item.name
                const type = item.typeline
                if (type.includes('Creature')){
                    image.setAttribute("data-trgtloc", "Creature")
                } else if (type.includes("Artifact")){
                    image.setAttribute("data-trgtloc", "Artifact")
                } else if (type.includes("Enchantment")){
                    image.setAttribute("data-trgtloc", "Enchantment")
                } else if (type.includes("Instant")){
                    image.setAttribute("data-trgtloc", "Instant")
                } else if (type.includes("Land")){
                    image.setAttribute("data-trgtloc", "Land")
                } else if (type.includes("Planeswalker")){
                    image.setAttribute("data-trgtloc", "Planeswalker")
                } else if (type.includes("Sorcery")){
                    image.setAttribute("data-trgtloc", "Sorcery")
                }
                lazyLoad(image)

                

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


    const di = Number(this.parentNode.getAttribute("data-index"))
    if (di == NaN){
        return
    }

    const oldCard = deckGrid.querySelector(`div.deck-card[data-index='${di}']`)

   



    //create the card holder
    const card = document.createElement("div")
    card.setAttribute("data-index", this.parentNode.getAttribute("data-index"))
    card.classList.add("deck-card")
   

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
    newimg.classList.add("deck-card-image")
    
    const loc = img.getAttribute("data-trgtloc")
    const segment = document.querySelector(`#${loc}`)


    card.append(del)
    card.append(commander)
    card.append(newimg)
    lazyLoad(card)
    segment.appendChild(card)

    
    
    deckAdd.push(di)
   
}

function deckClick(e){
    if (e.target.className === "delete-card"){
        const value = Number(e.target.parentNode.getAttribute("data-index"))

        if (value == NaN) {
            return 
        }

        const index = deckAdd.indexOf(value)
        if (index != -1){
            console.log(deckAdd)
            deckAdd.splice(index, 1)
            console.log(deckAdd)
        } else {
            deckRemove.push(value)
        }

        

        e.target.parentNode.remove()
    } else if (e.target.className === "make-commander"){
        const previousCom = document.querySelector("div.commander")
        if (previousCom != null){
            let comIndex = 0
            let temp = previousCom.getAttribute("data-index")
            comIndex = parseInt(temp)
            if (isNaN(comIndex)){
                comIndex = 0
            }
            
            deckAdd.push(comIndex)
            previousCom.classList.remove("commander")
        }
        
        const card = e.target.parentNode;
        

        if (card != null){
            let comIndex = 0
            card.classList.add("commander")
            let temp = card.getAttribute("data-index")
            comIndex = parseInt(temp)
            if (isNaN(comIndex)){
                console.log("NAN")
                return
            }

            const index = deckAdd.indexOf(comIndex)
            if (index != -1){
                deckAdd.splice(index, 1)
            }
        }
    }
    
}

async function saveDeck(e){
    const items = deckGrid.querySelectorAll('div.deck-card:not(.commander)')
    const deckName = document.querySelector('input.deck-name')

    
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

    let dID = deckGrid.getAttribute("data-index")
    if (dID == null){
        window.alert("Invalid DeckID")
    } 

    let dIDNum = Number(dID)
    if (dIDNum == NaN){
        return
    }



    const deckInfo = {
        'deckID': dIDNum,
        'commander': comIndex, 
        'addToDeck': deckAdd, 
        'removeFromDeck': deckRemove, 
        'deckName':deckName.value,
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
            window.alert("Bad Request")
        }

        resp = await response.json()
        deckGrid.setAttribute("data-index", resp.id)
    }catch(error){
        window.alert(error)
    }

}

function deleteClick(e){
    if (confirm("Are you sure you want to delete this deck?") == true){
        
        const id = deckGrid.getAttribute("data-index")
        try{
            fetch(`/cards/deleteDeck/${id}`, {
            method: 'Delete'})
            window.location.href ="https://localhost:4000/home"
        } catch(error){
            console.log(error)
            window.alert("Delete Unsuccsessful")
        }
    }
}

