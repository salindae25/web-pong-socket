import { useEffect, useRef } from "react";

export function GameBoard({
  height,
  width,
  state,
}: {
  height: number;
  width: number;
  state?: any;
}) {
  const canvasRef = useRef<any>(null);
  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    ctx.fillStyle = "gray";
    ctx.fillRect(0, 0, 400, 400);
  }, []);
  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    if (state) {
      ctx.clearRect(0, 0, 400, 400);
      ctx.fillStyle = "gray";
      ctx.fillRect(0, 0, 400, 400);
      ctx.fillStyle = "black";
      ctx.beginPath();
      ctx.arc(
        state?.Ball?.x,
        state?.Ball?.y,
        state?.Ball?.radius,
        0,
        2 * Math.PI
      );
      ctx.stroke();
      ctx.fillRect(state?.Player1?.x, state?.Player1?.y, 20, 100);
      ctx.fillRect(state?.Player2?.x, state?.Player2?.y, 20, 100);
    }
  }, [state]);
  return <canvas ref={canvasRef} id="canvas" height={height} width={width} />;
}
