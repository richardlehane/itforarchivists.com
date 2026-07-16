export async function onRequestGet(ctx) {
    const file = await ctx.env.SIEGFRIED.get(ctx.params.deb);
    if (!file) return new Response(null, { status: 404 });
    return new Response(file.body, {
        headers: { "Content-Type": file.httpMetadata.contentType },
    });
}