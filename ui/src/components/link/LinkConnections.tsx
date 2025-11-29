import { useLinkAnchor } from "@/lib/hooks/useLinkAnchor"

const makeBezierPath = (x1: number, y1: number, x2: number, y2: number) => {
  const dx = Math.abs(x2 - x1) * 0.3
  const cp1x = x1 + (x2 > x1 ? dx : -dx)
  const cp2x = x2 - (x2 > x1 ? dx : -dx)

  return `M ${x1},${y1} C ${cp1x},${y1} ${cp2x},${y2} ${x2},${y2}`
}

export const LinkConnections = () => {
  const { anchorsRef, connections, removeConnection, draggingFrom, tempPos, hoveredConnection, setHoveredConnection } = useLinkAnchor()

  return (
    <svg className="fixed top-0 left-0 w-screen h-screen pointer-events-none">
      <defs>
        <marker
          id="arrow"
          markerWidth="6"
          markerHeight="6"
          refX="5"
          refY="3"
          orient="auto"
        >
          <path d="M0,0 L6,3 L0,6 Z" fill="context-stroke" />
        </marker>
      </defs>

      {connections.map(({ from, to }) => {
        const a = anchorsRef.current[from]?.el?.getBoundingClientRect()
        const b = anchorsRef.current[to]?.el?.getBoundingClientRect()
        if (!a || !b) return null;

        const cid = from + "-" + to

        const x1 = a.left + a.width / 2
        const y1 = a.top + a.height / 2
        const x2 = b.left + b.width / 2
        const y2 = b.top + b.height / 2

        const path = makeBezierPath(x1, y1, x2, y2)
        const isHovered = hoveredConnection?.from === from && hoveredConnection?.to === to

        return (
          <g key={cid}>
            <path
              d={path}
              stroke="transparent"
              strokeWidth={16}
              fill="none"
              style={{ pointerEvents: "stroke" }}
              onPointerEnter={() => setHoveredConnection({ from, to })}
              onPointerLeave={() => setHoveredConnection(null)}
              onPointerDown={() => { removeConnection(from, to); setHoveredConnection(null) }}
              className="pointer-events-auto cursor-pointer"
            />

            <path
              d={path}
              stroke={isHovered ? "red" : "black"}
              strokeWidth={2}
              fill="none"
              markerEnd="url(#arrow)"
            />
          </g>
        )
      })}

      {draggingFrom && tempPos && (() => {
        const a = anchorsRef.current[draggingFrom].el?.getBoundingClientRect()
        if (!a) return null

        const x1 = a.left + a.width / 2
        const y1 = a.top + a.height / 2

        return (
          <path
            d={makeBezierPath(x1, y1, tempPos.x, tempPos.y)}
            stroke="gray"
            strokeDasharray="4"
            fill="none"
            strokeWidth={2}
          />
        )
      })()}
    </svg>
  )
}
