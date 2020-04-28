const record = document.querySelector('#record');
const shot = document.querySelector('#shot');
const hit = document.querySelector('#hit');
const dead = document.querySelector('#dead');
const myBoard = document.querySelector('#my-board');
const enemy = document.querySelector('#enemy');
const again = document.querySelector('#again');
const generate = document.querySelector('#generate');
const header = document.querySelector('.header');

let conn;

const boardSize = 10;
const horizontalDirectionProbability = 0.5;
const finalPhrase = 'Игра окончена!';
const finalColor = '#f00';
const storageRecordID = 'navalBattleRecord';
const myBoardId = 'm_';
const enemyBoardId = 'e_';

const game = {
  myShips: [],
  myShipCount: 0,
  ships: [],
  shipCount: 0,
  optionShip: {
    count: [1, 2, 3, 4],
    size: [4, 3, 2, 1],
  },
  collision: new Set(),
  generateShip(boardId) {
    switch (boardId) {
      case myBoardId:
        this.myShips = [];
        this.myShipCount = 0;
        break;
      case enemyBoardId:
        this.ships = [];
        this.shipCount = 0;
        break;
    }
    
    this.collision.clear();

    for (let i = 0; i < this.optionShip.count.length; i++) {
      for (let j = 0; j < this.optionShip.count[i]; j++) {
        const size = this.optionShip.size[i];
        const ship = this.generateOphionsShip(size, boardId);

        switch (boardId) {
          case myBoardId:
            this.myShips.push(ship);
            this.myShipCount++;
            break;
          case enemyBoardId:
            this.ships.push(ship);
            this.shipCount++;
            break;
        }
      }
    }
  },
  generateOphionsShip(shipSize, boardId) {
    const ship = {
      hit: [],
      location: [],
    };

    const direction = Math.random() < horizontalDirectionProbability;
    let x, y;

    if (direction) {
      x = Math.floor(Math.random() * boardSize);
      y = Math.floor(Math.random() * (boardSize - shipSize));
    } else {
      x = Math.floor(Math.random() * (boardSize - shipSize));
      y = Math.floor(Math.random() * boardSize);
    }

    for (let i = 0; i < shipSize; i++) {
      if (direction) {
        ship.location.push(boardId + x + '_' + (y + i));
      } else {
        ship.location.push(boardId + (x + i) + '_' + y);
      }
      ship.hit.push('');
    }

    if (this.checkCollision(ship.location)) {
      return this.generateOphionsShip(shipSize, boardId);
    }

    this.addCollision(ship.location, boardId);

    return ship;
  },
  checkCollision(location) {
    for (const coord of location) {
      if (this.collision.has(coord)) {
        return true;
      }
    }
    return false;
  },
  addCollision(location, boardId) {
    for (let i = 0; i < location.length; i++) {
      const coords = location[i].split('_');
      const startCoordX = coords[1] - 1;
      for (let x = startCoordX; x < startCoordX + 3; x++) {
        const startCoordY = coords[2] - 1;
        for (let y = startCoordY; y < startCoordY + 3; y++) {
          if (x >= 0 && x < boardSize && y >= 0 && y < boardSize) {
            const coord = boardId + x + '_' + y;
            this.collision.add(coord);
          }
        }
      }
    }
  }
};

const play = {
  record: localStorage.getItem(storageRecordID) || 0,
  shot: 0,
  hit: 0,
  dead: 0,

  /**
   * @param {string} data
   */
  set updateData(data) {
    this[data] += 1;
    this.render();
  },

  render() {
    record.textContent = this.record;
    shot.textContent = this.shot;
    hit.textContent = this.hit;
    dead.textContent = this.dead;
  },

  reset() {
    this.shot = 0;
    this.hit = 0;
    this.dead = 0;

    show.reset();

    game.generateShip(myBoardId);
    game.generateShip(enemyBoardId);

    this.render();
  }
};

const show = {
  hit(elem) {
    this.changeClass(elem, 'hit');
  },
  miss(elem) {
    this.changeClass(elem, 'miss');
  },
  dead(elem) {
    this.changeClass(elem, 'dead');
  },
  changeClass(elem, value) {
    elem.className = value;
  },
  showMyBoard() {
    for (const ship of game.myShips) {
      for (const location of ship.location) {
        myBoard.querySelector(`#${location}`).style.borderColor = '#f00';
      }
    }
  },
  reset() {
    enemy.querySelectorAll('td').forEach((cell) => {
      cell.classList.remove('hit', 'miss', 'dead');
    });

    myBoard.querySelectorAll('td').forEach((cell) => {
      cell.style.borderColor = '#000';
    });
  }
};

const generateBoard = (elem, prefix) => {
  const fragment = document.createDocumentFragment();
  for (let x = 0; x < boardSize; x++) {
    const tr = document.createElement('tr');
    for (let y = 0; y < boardSize; y++) {
      const td = document.createElement('td');
      td.id = prefix + x + '_' + y;
      tr.appendChild(td);
    }
    fragment.appendChild(tr);
  }
  elem.appendChild(fragment);
};

const connect = () => {
  let conn;

  if (window["WebSocket"]) {
    conn = new WebSocket(`ws://${document.location.host}/api/connect`);
    conn.onclose = (e) => {
      console.log('connection closed');
    };
    conn.onmessage = enemyStep;
  } else {
    console.log('You browser does not support WebSocket');
  }

  return conn;
};

const fire = (event) => {
  const target = event.target;

  if (target.classList.length > 0 || target.tagName !== 'TD' || game.shipCount < 1) {
    return;
  }

  conn.send(`step:${target.id}`);

  // show.miss(target);
  // play.updateData = 'shot';

  // for (let i = 0; i < game.ships.length; i++) {
  //   const ship = game.ships[i];
  //   const index = ship.location.indexOf(target.id);
  //   if (index >= 0) {
  //     show.hit(target);
  //     play.updateData = 'hit';
  //     ship.hit[index] = 'x';
  //     const life = ship.hit.indexOf('');
  //     if (life < 0) {
  //       play.updateData = 'dead';
  //       for (const id of ship.location) {
  //         show.dead(document.getElementById(id));
  //       }

  //       game.shipCount -= 1;

  //       if (game.shipCount < 1) {
  //         header.textContent = finalPhrase;
  //         header.style.color = finalColor;

  //         if (play.shot < play.record || play.record == 0) {
  //           localStorage.setItem(storageRecordID, play.shot);
  //           play.record = play.shot;
  //           play.render();
  //         }
  //       }
  //     }
  //   }
  // }
};

const enemyStep = (e) => {
  console.log(e.data);
  conn.send(`res:miss`);
};

const init = () => {
  generateBoard(myBoard, myBoardId);
  generateBoard(enemy, enemyBoardId);

  enemy.addEventListener('click', fire);
  play.render();

  again.addEventListener('click', () => {
    play.reset();
    conn = connect();
  });

  generate.addEventListener('click', () => {
    play.reset();
    game.generateShip(myBoardId);
    show.showMyBoard();
  });

  record.addEventListener('dblclick', () => {
    localStorage.clear();
    play.record = 0;
    play.render();
  });
};

document.addEventListener('DOMContentLoaded', init);
