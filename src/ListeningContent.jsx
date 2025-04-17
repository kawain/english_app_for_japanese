import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'

function ListeningContent ({
  level,
  onLevelChange,
  volume,
  onVolumeChange,
  isSoundEnabled,
  onToggleSound
}) {
  return (
    <>
      <h1>ヒアリング</h1>
      <VolumeControl
        volume={volume}
        onVolumeChange={onVolumeChange}
        isSoundEnabled={isSoundEnabled}
        onToggleSound={onToggleSound}
      />
      <LevelControl level={level} onLevelChange={onLevelChange} />
    </>
  )
}
export default ListeningContent
