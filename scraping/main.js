const { chromium } = require("playwright");
const Database = require("better-sqlite3");

const db = new Database("compras.db");

db.exec(`
CREATE TABLE IF NOT EXISTS comprador (
  id_comprador INTEGER PRIMARY KEY AUTOINCREMENT,
  nombre TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pedido (
  id_pedido INTEGER PRIMARY KEY AUTOINCREMENT,
  nombre TEXT NOT NULL,
  precio REAL NOT NULL,
  imagen TEXT,
  variante TEXT,
  cantidad INTEGER NOT NULL,
  id_comprador INTEGER,
  FOREIGN KEY (id_comprador) REFERENCES comprador(id_comprador)
);
`);

(async () => {
  const context = await chromium.launchPersistentContext("./perfil-temu", {
    headless: false,
    args: ["--disable-blink-features=AutomationControlled"],
  });

  const page = await context.newPage();

  try {
    await page.goto("https://www.temu.com/shopping_cart.html", {
      waitUntil: "domcontentloaded",
    });

    await page.waitForSelector("div._1is6WIRd", { timeout: 30000 });

    const productos = await page.evaluate(() => {
      const resultado = [];

      const productosDOM = document.querySelectorAll("div._1is6WIRd");

      productosDOM.forEach((prod) => {
        const nombreEl = prod.querySelector("a._13hQHbOY");
        const precioEl = prod.querySelector("._3FqV-jth span:last-child");
        const imgEl = prod.querySelector("img._2WPpsFlf");
        const varianteEl = prod.querySelector("div._24c2Me98");
        const cantidadEl = prod.querySelector("input._3rRk6Q66");

        if (!nombreEl || !precioEl) return;

        const nombre = nombreEl.innerText.trim();

        const precioTexto = precioEl.innerText
          .replace(",", ".")
          .replace(/[^0-9.]/g, "");

        const precio = parseFloat(precioTexto);

        const imagen = imgEl ? imgEl.src : null;

        const variante = varianteEl ? varianteEl.innerText.trim() : null;

        const cantidad = cantidadEl ? parseInt(cantidadEl.value) : 1;

        resultado.push({
          nombre,
          precio,
          imagen,
          variante,
          cantidad,
          id_comprador: null,
        });
      });

      return resultado;
    });

    const insert = db.prepare(`
      INSERT INTO pedido (nombre, precio, imagen, variante, cantidad, id_comprador)
      VALUES (@nombre, @precio, @imagen, @variante, @cantidad, @id_comprador)
    `);

    const insertMany = db.transaction((items) => {
      for (const item of items) insert.run(item);
    });

    insertMany(productos);

    console.log(`✅ Guardados ${productos.length} productos en compras.db`);
  } catch (error) {
    console.error("❌ Error:", error);
  } finally {
    await context.close();
    db.close();
  }
})();
