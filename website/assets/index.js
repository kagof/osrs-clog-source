const sources = await getSources()

let lastUpdated = document.getElementById('last-updated')
lastUpdated.innerText = lastUpdated.innerText.replace('unknown', new Date(sources.lastUpdated).toLocaleString("en-GB", {
  day: "2-digit",
  month: "short",
  year: "numeric",
  hour: "numeric",
  minute: "2-digit",
  hour12: true,
  timeZoneName: "short",
}))

let searchElement = document.getElementById('searchbar-input')
searchElement.onkeydown = (event) => {
    if (event.keyCode == 27) { // esc
        clearDropdown()
        event.srcElement.value = ''
    }
    if (event.keyCode !== 13) { // enter/return
        return
    }
    let value = event.srcElement.value
    if (value == null) {
        return
    }
    console.log(`searching for ${value}`)
    let results = Object.keys(sources.sources)
        .filter((name) => name.toLowerCase()
            .split(' ')
            .some((word) => word.startsWith(value.toLowerCase())))
    console.log(`${results.length} results found`)
    setDropdownItems(results)
}

async function getSources() {
    let storedHash = window.localStorage.getItem('sources.json.hash')
    if (!storedHash) {
        return fetchSourcesFromRemote()
    }
    let storageVal = window.localStorage.getItem('sources.json')
    if (storageVal != null) {
        console.log('sources found in local storage')
        let remoteHash = await (await fetch('data/sources.json.hash')).text()
        if (remoteHash != storedHash) {
            console.log('hash mismatch; fetching sources')
            return fetchSourcesFromRemote()
        }
        try {
            console.log('hashes match, using local storage sources')
            return parseJsonAsync(storageVal)
        } catch (e) {
            console.log(`error parsing sources from local storage: ${e}`)
            return fetchSourcesFromRemote()
        }
    }
    console.log('no local storage of sources, fetching')
    return fetchSourcesFromRemote()
}

async function fetchSourcesFromRemote() {
    let response = await fetch('data/sources.json')
    let textData = await response.text()
    let hash = await sha256(textData)
    window.localStorage.setItem('sources.json', textData)
    window.localStorage.setItem('sources.json.hash', hash)
    return parseJsonAsync(textData)
}

async function sha256(str) {
  const data = new TextEncoder().encode(str);
  const hash = await crypto.subtle.digest("SHA-256", data);

  return [...new Uint8Array(hash)]
    .map(b => b.toString(16).padStart(2, "0"))
    .join("");
}

async function parseJsonAsync(str) {
    return new Promise((resolve, reject) => {
            try {
                resolve(JSON.parse(str))
            } catch (e) {
                reject(e)
            }
        })
}

let dropdownMenu = document.getElementById('dropdown-menu')
let dropdownItem0 = document.getElementById('dropdown-item-0')
let dropdownItem1 = document.getElementById('dropdown-item-1')
let dropdownItem2 = document.getElementById('dropdown-item-2')
let dropdownItem3 = document.getElementById('dropdown-item-3')
let dropdownItem4 = document.getElementById('dropdown-item-4')
let dropdownItemMore = document.getElementById('dropdown-item-more')
let dropdownItemEntry0 = document.getElementById('dropdown-item-entry-0')
let dropdownItemEntry1 = document.getElementById('dropdown-item-entry-1')
let dropdownItemEntry2 = document.getElementById('dropdown-item-entry-2')
let dropdownItemEntry3 = document.getElementById('dropdown-item-entry-3')
let dropdownItemEntry4 = document.getElementById('dropdown-item-entry-4')
let dropdownItemEntryMore = document.getElementById('dropdown-item-entry-more')

async function setDropdownItems(results) {
    clearDropdown()
    dropdownMenu.hidden = false
    if (results == null || results.length === 0) {
        dropdownItemEntry0.textContent = 'No results'
        dropdownItem0.hidden = false
        return
    }
    dropdownItemEntry0.textContent = results[0]
    dropdownItem0.hidden = false
    if (results.length < 2) {
        return
    }
    dropdownItemEntry1.textContent = results[1]
    dropdownItem1.hidden = false
    if (results.length < 3) {
        return
    }
    dropdownItemEntry2.textContent = results[2]
    dropdownItem2.hidden = false
    if (results.length < 4) {
        return
    }
    dropdownItemEntry3.textContent = results[3]
    dropdownItem3.hidden = false
    if (results.length < 5) {
        return
    }
    dropdownItemEntry4.textContent = results[4]
    dropdownItem4.hidden = false
    if (results.length === 5) {
        return
    }
    dropdownItemEntryMore.textContent = `and ${results.length - 5} more...`
    dropdownItemMore.hidden = false
}

async function clearDropdown() {
    dropdownMenu.hidden = true
    dropdownItem0.hidden = true
    dropdownItem1.hidden = true
    dropdownItem2.hidden = true
    dropdownItem3.hidden = true
    dropdownItem4.hidden = true
    dropdownItemMore.hidden = true
    dropdownItemEntry0.textContent = ''
    dropdownItemEntry1.textContent = ''
    dropdownItemEntry2.textContent = ''
    dropdownItemEntry3.textContent = ''
    dropdownItemEntry4.textContent = ''
    dropdownItemEntryMore.textContent = ''
}

let table = document.getElementById('search-results-table')
let searchResults = document.getElementById('search-results')

function clearChoice() {
    sessionStorage.removeItem('choice')
    searchResults.hidden = true
    while (table.rows.length > 0) {
        table.deleteRow(-1)
    }
}

function insertSourceRow(sourceName) {
    let row = table.insertRow()
    let th = document.createElement('th');
    th.colSpan = 5
    th.className = 'search-results-table-source'
    th.textContent = sourceName;
    row.appendChild(th);
}

function insertSubclassificationRow(subclassification) {
    let row = table.insertRow()
    let th = document.createElement('th')
    th.colSpan = 5
    th.className = 'search-results-table-subclassification'
    let a = document.createElement('a')
    a.href = subclassification.link
    a.textContent = !subclassification.name ? 'Standard' : subclassification.name
    th.appendChild(a)
    row.appendChild(th)
}

function insertItemHeadersRow() {
    let row = table.insertRow()
    row.className = 'search-results-table-items-headers'
    
    let th0 = document.createElement('th')
    th0.textContent = 'Item'
    th0.colSpan = 2
    let th1 = document.createElement('th')
    th1.textContent = 'Rarity'
    let th2 = document.createElement('th')
    th2.textContent = 'Quantity'
    let th3 = document.createElement('th')
    th3.textContent = 'Percent obtained'
    row.appendChild(th0)
    row.appendChild(th1)
    row.appendChild(th2)
    row.appendChild(th3)
}

function insertItemRow(item) {
    let row = table.insertRow()
    row.className = 'search-results-table-item'

    let cell0 = row.insertCell()
    let img = document.createElement('img')
    img.src = item.image
    cell0.appendChild(img)
    
    let cell1 = row.insertCell()
    let a = document.createElement('a')
    a.href = item.link
    a.textContent = item.name
    cell1.appendChild(a)
    
    let cell2 = row.insertCell()
    cell2.textContent = item.rarity
    
    let cell3 = row.insertCell()
    cell3.textContent = item.quantity

    let cell4 = row.insertCell()
    cell4.textContent = item.compPercent
}

let clearButton = document.getElementById('clear-button')
clearButton.onclick = clearChoice

function dropdownItemEntryOnclick(event) {
    let key = event.srcElement.textContent
    setChoice(key)
}

function setChoice(key) {
    console.log(`${key} selected`)
    let source = sources.sources[key]
    if (!source) {
        console.log(`unknown choice ${key}`)
        clearChoice()
        return
    }
    clearDropdown()
    clearChoice()
    insertSourceRow(key)
    Object.values(source.subclassifications)
    // TODO: sort
    .forEach(subclassification => {
        insertSubclassificationRow(subclassification)
        insertItemHeadersRow()
        subclassification.items.forEach(insertItemRow) // TODO: sort
    })
    searchResults.hidden = false
    sessionStorage.setItem('choice', key)
}

let initChoice = sessionStorage.getItem('choice')
if (!!initChoice) {
    setChoice(initChoice)
}

dropdownItemEntry0.onclick = dropdownItemEntryOnclick
dropdownItemEntry1.onclick = dropdownItemEntryOnclick
dropdownItemEntry2.onclick = dropdownItemEntryOnclick
dropdownItemEntry3.onclick = dropdownItemEntryOnclick
dropdownItemEntry4.onclick = dropdownItemEntryOnclick
