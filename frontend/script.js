const record = document.getElementById('record');
const shot = document.getElementById('shot');
const hit = document.getElementById('hit');
const dead = document.getElementById('dead');
const enemy = document.getElementById('enemy');
const again = document.getElementById('again');
const header = document.querySelector('.header');

const boardSize = 10;
const horizontalDirectionProbability = 0.5;
const finalPhrase = 'Игра окончена!';
const finalColor = '#f00';
const storageRecordID = 'navalBattleRecord';

const game = {
  ships: [],
  shipCount: 0,
  optionShip: {
    count: [1, 2, 3, 4],
    size: [4, 3, 2, 1],
  },
  collision: new Set(),
  generateShip() {
    this.ships = [];
    this.shipCount = 0;
    this.collision.clear();

    for (let i = 0; i < this.optionShip.count.length; i++) {
      for (let j = 0; j < this.optionShip.count[i]; j++) {
        const size = this.optionShip.size[i];
        const ship = this.generateOphionsShip(size);
        this.ships.push(ship);
        this.shipCount++;
      }
    }
  },
  generateOphionsShip(shipSize) {
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
        ship.location.push(x + '' + (y + i));
      } else {
        ship.location.push((x + i) + '' + y);
      }
      ship.hit.push('');
    }

    if (this.checkCollision(ship.location)) {
      return this.generateOphionsShip(shipSize);
    }

    this.addCollision(ship.location);

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
  addCollision(location) {
    for (let i = 0; i < location.length; i++) {
      const startCoordX = location[i][0] - 1;
      for (let x = startCoordX; x < startCoordX + 3; x++) {
        const startCoordY = location[i][1] - 1;
        for (let y = startCoordY; y < startCoordY + 3; y++) {
          if (x >= 0 && x < boardSize && y >= 0 && y < boardSize) {
            const coord = x + '' + y;
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

    enemy.querySelectorAll('td').forEach((cell) => {
      cell.classList.remove('hit', 'miss', 'dead');
    });

    game.generateShip();

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
};

const fire = (event) => {
  const target = event.target;

  if (target.classList.length > 0 || target.tagName !== 'TD' || game.shipCount < 1) {
    return;
  }

  show.miss(target);
  play.updateData = 'shot';

  for (let i = 0; i < game.ships.length; i++) {
    const ship = game.ships[i];
    const index = ship.location.indexOf(target.id);
    if (index >= 0) {
      show.hit(target);
      play.updateData = 'hit';
      ship.hit[index] = 'x';
      const life = ship.hit.indexOf('');
      if (life < 0) {
        play.updateData = 'dead';
        for (const id of ship.location) {
          show.dead(document.getElementById(id));
        }

        game.shipCount -= 1;

        if (game.shipCount < 1) {
          header.textContent = finalPhrase;
          header.style.color = finalColor;

          if (play.shot < play.record || play.record == 0) {
            localStorage.setItem(storageRecordID, play.shot);
            play.record = play.shot;
            play.render();
          }
        }
      }
    }
  }
};

const init = () => {
  enemy.addEventListener('click', fire);
  play.render();
  game.generateShip();

  again.addEventListener('click', () => {
    play.reset();
  });

  record.addEventListener('dblclick', () => {
    localStorage.clear();
    play.record = 0;
    play.render();
  });
};

init();
