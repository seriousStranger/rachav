import { QRCode } from "react-qr-code"
import { Button } from "@/components/ui/button"
import { RiCloseLine } from "@remixicon/react"

interface QRCodePopupProps {
  link: string
  onClose: () => void
}

export function QRCodePopup({ link, onClose }: QRCodePopupProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-background rounded-lg border shadow-lg max-w-sm w-full p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">QR Code</h3>
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={onClose}
            aria-label="Close"
          >
            <RiCloseLine className="h-4 w-4" />
          </Button>
        </div>
        <div className="flex flex-col items-center gap-4">
          <div className="bg-white p-4 rounded border">
            <QRCode value={link} size={200} />
          </div>
          <p className="text-sm text-muted-foreground text-center break-all">
            {link}
          </p>
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigator.clipboard.writeText(link)}
          >
            Copy Link
          </Button>
        </div>
      </div>
    </div>
  )
}