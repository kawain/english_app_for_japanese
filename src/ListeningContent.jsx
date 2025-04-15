import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'

function ListeningContent ({ volume, onVolumeChange, level, onLevelChange }) {
  return (
    <>
      <h1>ヒアリング</h1>
      <VolumeControl volume={volume} onVolumeChange={onVolumeChange} />
      <LevelControl level={level} onLevelChange={onLevelChange} />
    </>
  )
}
export default ListeningContent
