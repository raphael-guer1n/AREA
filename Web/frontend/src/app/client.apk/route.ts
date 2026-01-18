import { createReadStream } from "fs";
import { stat } from "fs/promises";
import { NextResponse } from "next/server";
import { Readable } from "stream";

const DEFAULT_APK_PATH = "/apk/client.apk";

export async function GET(): Promise<NextResponse> {
  const apkPath = process.env.APK_FILE_PATH?.trim() || DEFAULT_APK_PATH;

  try {
    const fileStat = await stat(apkPath);
    if (!fileStat.isFile()) {
      return NextResponse.json({ error: "APK introuvable." }, { status: 404 });
    }

    const stream = Readable.toWeb(createReadStream(apkPath)) as ReadableStream;
    return new NextResponse(stream, {
      status: 200,
      headers: {
        "Content-Type": "application/vnd.android.package-archive",
        "Content-Length": fileStat.size.toString(),
        "Content-Disposition": 'attachment; filename="client.apk"',
      },
    });
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Impossible de récupérer l'APK.";
    return NextResponse.json({ error: message }, { status: 404 });
  }
}
