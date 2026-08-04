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
    await page.goto("https://www.temu.com/cart", {
      waitUntil: "networkidle",
    });

    await page.waitForSelector("div.g_p_gl", { timeout: 30000 });

    const productos = await page.evaluate(() => {
      const resultado = [];
      const productosDOM = document.querySelectorAll("div._1is6WIRd");

      productosDOM.forEach((prod) => {
        const nombreEl = prod.querySelector("a._2Tl9qLr1");
        const precioEl = prod.querySelector("span._2hmmv2gR");
        const precio = precioEl
          ? parseFloat(
              precioEl.innerText
                .replace(/[^0-9,.]/g, "")
                .replace(",", ".")
            )
          : null;

        if (nombreEl && precio !== null && !Number.isNaN(precio)) {
          const nombre = nombreEl.innerText.trim();

          const imgEl = prod.querySelector("div._2PQb6lDA img.wxWpAMbp");
          const imagen = imgEl ? imgEl.src : null;

          const varianteEl = prod.querySelector("div._1TbbNTn_ span");
          const variante = varianteEl ? varianteEl.innerText.trim() : "N/A";

          const cantidadEl = prod.querySelector("input[aria-label]");
          const cantidad = cantidadEl ? parseInt(cantidadEl.value) : 1;

          resultado.push({
            nombre,
            precio,
            imagen,
            variante,
            cantidad,
            id_comprador: null,
          });
        }
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

    console.log(
      `✅ Se han guardado ${productos.length} productos en la base de datos compras.db`
    );
  } catch (error) {
    console.error("❌ Error durante el proceso:", error);
  } finally {
    await context.close();
    db.close();
  }
})();
