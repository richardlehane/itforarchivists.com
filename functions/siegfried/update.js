export async function onRequestGet(ctx) {
    const url = new URL("/siegfried/update/update.json");
    return await ctx.env.ASSETS.fetch(url);
}