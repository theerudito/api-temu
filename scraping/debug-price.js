const { chromium } = require("playwright");
const fs = require("fs");

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

    const debugInfo = await page.evaluate(() => {
      const productosDOM = document.querySelectorAll("div._1is6WIRd");
      if (productosDOM.length === 0) return { error: "No se encontraron productos" };

      // Buscar el producto que tenga el texto "después de las promociones"
      const prodConPromo = Array.from(productosDOM).find((prod) =>
        prod.innerText.includes("después de las promociones")
      );

      if (!prodConPromo) return { error: "no encontrado" };

      return { html: prodConPromo.outerHTML };
    });

    fs.writeFileSync("debug-output.txt", JSON.stringify(debugInfo, null, 2));
    console.log("listo");
  } catch (error) {
    console.error("❌ Error durante el proceso:", error);
  } finally {
    await context.close();
  }
})();
